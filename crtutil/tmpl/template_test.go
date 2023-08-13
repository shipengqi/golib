package tmpl

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shipengqi/golib/crtutil"
)

func TestBuildDefaultCertTemplate(t *testing.T) {
	tests := []struct {
		title    string
		input    string
		expected []string
	}{
		{
			"successfully output the certificate according to the given template",
			"../testdata/server.crt",
			[]string{
				"Serial: 4751997750760398084",
				"Valid: 2021-11-29 08:39 UTC to 2022-11-29 08:39 UTC",
				"Signature: SHA256-RSA",
				"BitLength: 2048",
				"Authority Key ID: CA:A5:79:D4:EB:5D:1F:F0:8F:40:52:A9:AF:3B:E7:6B:84:74:F9:B9",
				"Basic Constraints: CA:false, pathlen:-1",
			},
		},
		{
			"successfully output the self-signed certificate according to the given template",
			"../testdata/self-signed.crt",
			[]string{
				"Serial: 5577006791947779410",
				"Valid: 2022-09-23 06:09 UTC to 2032-09-30 06:09 UTC",
				"Signature: SHA256-RSA (self-signed)",
				"BitLength: 4096",
				"Subject Key ID: 6D:E9:2B:2B:1D:59:AB:B5:46:8C:7B:93:C3:49:7E:95:B0:20:E5:4C",
				"Basic Constraints: CA:true, pathlen:-1",
			},
		},
	}

	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			parsed, err := crtutil.ReadAsX509FromFile(v.input)
			assert.NoError(t, err)
			assert.NotEmpty(t, parsed)

			output, err := BuildDefaultCertTemplate(parsed[0], true)
			assert.NoError(t, err)

			for _, contain := range v.expected {
				assert.Contains(t, string(output), contain)
			}
		})
	}
}
