package tmpl

import (
	"bytes"
	"crypto/x509"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

const (
	VerboseTmpl = `
{{- define "PkixName" -}}
{{- range .Names}}
	{{ .Type | oidName }}: {{ .Value }}
{{- end -}}
{{end -}}

{{- define "IsSelfSigned" -}}
{{ . | isSelfSigned }}
{{end -}}

Serial: {{.SerialNumber}}
Valid: {{.NotBefore | notBefore}} to {{.NotAfter | notAfter}}
Signature: {{.SignatureAlgorithm | highlightAlgorithm}}{{- template "IsSelfSigned" . -}}
{{- if (showBitLen .) }}
BitLength: {{.PublicKey | bitLen}}
{{end}}
Subject Info:
	{{- template "PkixName" .Subject}}
Issuer Info:
	{{- template "PkixName" .Issuer}}
{{- if .SubjectKeyId}}
Subject Key ID: {{.SubjectKeyId | tohex}}
{{- end}}
{{- if .AuthorityKeyId}}
Authority Key ID: {{.AuthorityKeyId | tohex}}
{{- end}}
{{- if .BasicConstraintsValid}}
Basic Constraints: CA:{{.IsCA}}{{if .MaxPathLen}}, pathlen:{{.MaxPathLen}}{{end}}{{end}}
{{- if (nameConstraints .) }}
Name Constraints{{if .PermittedDNSDomainsCritical}} (critical){{end}}: 
{{- if .PermittedDNSDomains}}
Permitted DNS domains:
	{{join ", " .PermittedDNSDomains}}
{{- end -}}
{{- if .PermittedEmailAddresses}}
Permitted email addresses:
	{{join ", " .PermittedEmailAddresses}}
{{- end -}}
{{- if .PermittedIPRanges}}
Permitted IP ranges:
	{{join ", " .PermittedIPRanges}}
{{- end -}}
{{- if .PermittedURIDomains}}
Permitted URI domains:
	{{join ", " .PermittedURIDomains}}
{{- end}}
{{- if .ExcludedDNSDomains}}
Excluded DNS domains:
	{{join ", " .ExcludedDNSDomains}}
{{- end}}
{{- if .ExcludedEmailAddresses}}
Excluded email addresses:
	{{join ", " .ExcludedEmailAddresses}}
{{- end}}
{{- if .ExcludedIPRanges}}
Excluded IP ranges:
	{{join ", " .ExcludedIPRanges}}
{{- end}}
{{- if .ExcludedURIDomains}}
Excluded URI domains:
	{{join ", " .NameConstraints.ExcludedURIDomains}}
{{- end}}
{{- end}}
{{- if .OCSPServer}}
OCSP Server(s):
	{{join ", " .OCSPServer}}
{{- end}}
{{- if .IssuingCertificateURL}}
Issuing Certificate URL(s):
	{{join ", " .IssuingCertificateURL}}
{{- end}}
{{- if .KeyUsage}}
Key Usage:
{{- range .KeyUsage | keyUsage}}
	{{.}}
{{- end}}
{{- end}}
{{- if .ExtKeyUsage}}
Extended Key Usage:
{{- range .ExtKeyUsage}}
	{{. | extKeyUsage}}{{end}}
{{- end}}
{{- if .DNSNames}}
DNS Names:
	{{join ", " .DNSNames}}
{{- end}}
{{- if .IPAddresses}}
IP Addresses:
	{{join ", " .IPAddresses}}
{{- end}}
{{- if .URIs}}
URI Names:
	{{join ", " .URIs}}
{{- end}}
{{- if .EmailAddresses}}
Email Addresses:
	{{join ", " .EmailAddresses}}
{{- end}}`

	SimpleTmpl = `Valid: {{.NotBefore | notBefore}} to {{.NotAfter | notAfter}}
Subject:
	{{.Subject | shortName}}
Issuer:
	{{.Issuer | shortName}}
{{- if .DNSNames}}
DNS Names:
	{{join ", " .DNSNames}}{{end}}
{{- if .IPAddresses}}
IP Addresses:
	{{join ", " .IPAddresses}}{{end}}
{{- if .URIs}}
URI Names:
	{{join ", " .URIs}}{{end}}
{{- if .EmailAddresses}}
Email Addresses:
	{{join ", " .EmailAddresses}}{{end}}`

	WarningTmpl = `
{{- if .Warnings}}
Warnings:
{{- range .Warnings}}
	{{colorize . "red"}}
{{- end}}
{{- end}}`
)

// BuildCertFuncMap build a template.FuncMap with some extras.
func BuildCertFuncMap() template.FuncMap {
	funcmap := sprig.TxtFuncMap()
	extras := template.FuncMap{
		"notBefore":          NotBefore,
		"notAfter":           NotAfter,
		"colorize":           Colorize,
		"highlightAlgorithm": HighlightAlgorithm,
		"tohex":              Hexadecimalize,
		"keyUsage":           KeyUsage,
		"extKeyUsage":        ExtKeyUsage,
		"oidName":            OidName,
		"oidShort":           OidShort,
		"shortName":          ShortName,
		"commonName":         CommonName,
		"isSelfSigned":       ShowSelfSigned,
		"nameConstraints":    ShowNameConstraints,
		"showBitLen":         ShowBitLength,
		"bitLen":             BitLength,
	}
	for k, v := range extras {
		funcmap[k] = v
	}
	return funcmap
}

func BuildDefaultCertTemplate(cert *x509.Certificate, verbose bool) ([]byte, error) {
	t := template.New("crtutil template").Funcs(BuildCertFuncMap())

	var err error
	if verbose {
		t, err = t.Parse(VerboseTmpl)
	} else {
		t, err = t.Parse(SimpleTmpl)
	}
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, cert)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
