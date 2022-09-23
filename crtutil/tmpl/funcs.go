package tmpl

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/shipengqi/golib/crtutil"
)

const (
	ValidityWarnThreshold = "720h"
	TimeFormat            = "2006-01-02 15:04 MST"
)

var (
	Green  = color.New(color.Bold, color.FgGreen)
	Yellow = color.New(color.Bold, color.FgYellow)
	Red    = color.New(color.Bold, color.FgRed)
	Blue   = color.New(color.Bold, color.FgBlue)
	Cyan   = color.New(color.Bold, color.FgCyan)
)

var AlgorithmMapping = map[x509.SignatureAlgorithm]AlgorithmDesc{
	x509.MD2WithRSA:      {Red, "MD2-RSA"},
	x509.MD5WithRSA:      {Red, "MD5-RSA"},
	x509.SHA1WithRSA:     {Red, "SHA1-RSA"},
	x509.SHA256WithRSA:   {Green, "SHA256-RSA"},
	x509.SHA384WithRSA:   {Green, "SHA384-RSA"},
	x509.SHA512WithRSA:   {Green, "SHA512-RSA"},
	x509.DSAWithSHA1:     {Red, "DSA-SHA1"},
	x509.DSAWithSHA256:   {Red, "DSA-SHA256"},
	x509.ECDSAWithSHA1:   {Red, "ECDSA-SHA1"},
	x509.ECDSAWithSHA256: {Green, "ECDSA-SHA256"},
	x509.ECDSAWithSHA384: {Green, "ECDSA-SHA384"},
	x509.ECDSAWithSHA512: {Green, "ECDSA-SHA512"},
	x509.PureEd25519:     {Green, "ED25519"},
}

var KeyUsageStringMapping = map[x509.KeyUsage]string{
	x509.KeyUsageDigitalSignature:  "Digital Signature",
	x509.KeyUsageContentCommitment: "Content Commitment",
	x509.KeyUsageKeyEncipherment:   "Key Encipherment",
	x509.KeyUsageDataEncipherment:  "Data Encipherment",
	x509.KeyUsageKeyAgreement:      "Key Agreement",
	x509.KeyUsageCertSign:          "Cert Sign",
	x509.KeyUsageCRLSign:           "CRL Sign",
	x509.KeyUsageEncipherOnly:      "Encipher Only",
	x509.KeyUsageDecipherOnly:      "Decipher Only",
}

var ExtKeyUsageStringMapping = map[x509.ExtKeyUsage]string{
	x509.ExtKeyUsageAny:                        "Any",
	x509.ExtKeyUsageServerAuth:                 "Server Auth",
	x509.ExtKeyUsageClientAuth:                 "Client Auth",
	x509.ExtKeyUsageCodeSigning:                "Code Signing",
	x509.ExtKeyUsageEmailProtection:            "Email Protection",
	x509.ExtKeyUsageIPSECEndSystem:             "IPSEC End System",
	x509.ExtKeyUsageIPSECTunnel:                "IPSEC Tunnel",
	x509.ExtKeyUsageIPSECUser:                  "IPSEC User",
	x509.ExtKeyUsageTimeStamping:               "Time Stamping",
	x509.ExtKeyUsageOCSPSigning:                "OCSP Signing",
	x509.ExtKeyUsageMicrosoftServerGatedCrypto: "Microsoft ServerGatedCrypto",
	x509.ExtKeyUsageNetscapeServerGatedCrypto:  "Netscape ServerGatedCrypto",
}

var KeyUsages = []x509.KeyUsage{
	x509.KeyUsageDigitalSignature,
	x509.KeyUsageContentCommitment,
	x509.KeyUsageKeyEncipherment,
	x509.KeyUsageDataEncipherment,
	x509.KeyUsageKeyAgreement,
	x509.KeyUsageCertSign,
	x509.KeyUsageCRLSign,
	x509.KeyUsageEncipherOnly,
	x509.KeyUsageDecipherOnly,
}

// OidDesc returns a human-readable name, a short acronym from RFC1485, a snake_case slug suitable as a json key,
// and a boolean describing whether multiple copies can appear on an X509 cert.
type OidDesc struct {
	Name     string
	Short    string
	Slug     string
	Multiple bool
}

type AlgorithmDesc struct {
	Color *color.Color
	Name  string
}

func (a *AlgorithmDesc) String() string {
	return a.Color.SprintFunc()(a.Name)
}

func CommonName(name pkix.Name) string {
	if name.CommonName != "" {
		return fmt.Sprintf("CN=%s", name.CommonName)
	}
	return ShortName(name)
}

func ShortName(name pkix.Name) (out string) {
	show := false
	for _, n := range name.Names {
		short := OidShort(n.Type)
		if short != "" {
			if show {
				out += ", "
			}
			out += fmt.Sprintf("%s=%v", short, n.Value)
			show = true
		}
	}

	return
}

// HighlightAlgorithm changes the color of the signing algorithm.
func HighlightAlgorithm(alg x509.SignatureAlgorithm) string {
	desc, ok := AlgorithmMapping[alg]
	if !ok {
		return alg.String()
	}
	return desc.String()
}

// NotBefore takes a given NotBefore time of a certificate
// and returns that colorized time properly based on how
func NotBefore(start time.Time) string {
	return start.Format(TimeFormat)
}

// NotAfter takes a given NotAfter time of a certificate
// and returns that colorized time properly based on how
// close it is to expiry.
// If the certificate is valid for more than one month
// returns a green string.
// If the certificate is valid is less than one month
// the string will be yellow.
// If the certificate has already expired, the string
// will be red.
func NotAfter(end time.Time) string {
	now := time.Now()
	threshold := thresholdToTime(ValidityWarnThreshold, now)

	if end.Before(now) {
		return ColorizeTimeString(end, "red")
	} else if end.Before(threshold) {
		return ColorizeTimeString(end, "yellow")
	} else {
		return ColorizeTimeString(end, "green")
	}
}

// KeyUsage returns key usage string from a certificate.
func KeyUsage(ku x509.KeyUsage) []string {
	var out []string
	for _, key := range KeyUsages {
		if ku&key > 0 {
			out = append(out, KeyUsageStringMapping[key])
		}
	}
	return out
}

// ExtKeyUsage returns extended key usage string from a certificate.
func ExtKeyUsage(eku x509.ExtKeyUsage) string {
	val, ok := ExtKeyUsageStringMapping[eku]
	if ok {
		return val
	}
	return fmt.Sprintf("Unknown:%d", eku)
}

// Colorize colorizes the given string.
func Colorize(text, c string) string {
	switch strings.ToLower(c) {
	case "red":
		return Red.SprintfFunc()(text)
	case "yellow":
		return Yellow.SprintfFunc()(text)
	case "green":
		return Green.SprintfFunc()(text)
	case "blue":
		return Blue.SprintfFunc()(text)
	case "cyan":
		return Cyan.SprintfFunc()(text)
	default:
		return Blue.SprintfFunc()(text)
	}
}

// Hexadecimalize returns a colon separated, hexadecimal representation
// of a given byte array.
func Hexadecimalize(data []byte) string {
	var hexed bytes.Buffer
	for i := 0; i < len(data); i++ {
		hexed.WriteString(strings.ToUpper(hex.EncodeToString(data[i : i+1])))
		if i < len(data)-1 {
			hexed.WriteString(":")
		}
	}
	return hexed.String()
}

func OidName(oid asn1.ObjectIdentifier) string {
	return oiddesc(oid).Name
}

func OidShort(oid asn1.ObjectIdentifier) string {
	return oiddesc(oid).Short
}

func ColorizeTimeString(t time.Time, c string) string {
	return Colorize(t.Format(TimeFormat), c)
}

func ShowNameConstraints(cert *x509.Certificate) bool {
	if cert.PermittedDNSDomains != nil || cert.PermittedEmailAddresses != nil ||
		cert.PermittedIPRanges != nil || cert.PermittedURIDomains != nil ||
		cert.ExcludedDNSDomains != nil || cert.ExcludedEmailAddresses != nil ||
		cert.ExcludedIPRanges != nil || cert.ExcludedURIDomains != nil {
		return true
	}

	return false
}

func ShowSelfSigned(cert *x509.Certificate) string {
	if crtutil.IsSelfSigned(cert) {
		return " (self-signed)"
	}

	return ""
}

func thresholdToTime(threshold string, nowT ...time.Time) time.Time {
	var now time.Time
	if len(nowT) == 0 {
		now = time.Now()
	} else {
		now = nowT[0]
	}
	month, _ := time.ParseDuration(threshold)
	return now.Add(month)
}

func oiddesc(oid asn1.ObjectIdentifier) OidDesc {
	raw := oid.String()
	// Multiple should be true for any types that are []string in x509.pkix.Name. When in doubt, set it to true.
	names := map[string]OidDesc{
		"2.5.4.3":                    {"CommonName", "CN", "common_name", false},
		"2.5.4.5":                    {"EV Incorporation Registration Number", "", "ev_registration_number", false},
		"2.5.4.6":                    {"Country", "C", "country", true},
		"2.5.4.7":                    {"Locality", "L", "locality", true},
		"2.5.4.8":                    {"Province", "ST", "province", true},
		"2.5.4.9":                    {"Street", "", "street", true},
		"2.5.4.10":                   {"Organization", "O", "organization", true},
		"2.5.4.11":                   {"Organizational Unit", "OU", "organizational_unit", true},
		"2.5.4.15":                   {"Business Category", "", "business_category", true},
		"2.5.4.17":                   {"Postal Code", "", "postalcode", true},
		"1.2.840.113549.1.9.1":       {"Email Address", "", "email_address", true},
		"1.3.6.1.4.1.311.60.2.1.1":   {"EV Incorporation Locality", "", "ev_locality", true},
		"1.3.6.1.4.1.311.60.2.1.2":   {"EV Incorporation Province", "", "ev_province", true},
		"1.3.6.1.4.1.311.60.2.1.3":   {"EV Incorporation Country", "", "ev_country", true},
		"0.9.2342.19200300.100.1.1":  {"User ID", "UID", "user_id", true},
		"0.9.2342.19200300.100.1.25": {"Domain Component", "DC", "domain_component", true},
	}
	if desc, ok := names[raw]; ok {
		return desc
	}
	return OidDesc{raw, "", raw, true}
}
