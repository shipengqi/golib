package cliutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieveFlag(t *testing.T) {
	type res struct {
		value string
		ok    bool
	}
	tests := []struct {
		long     string
		short    string
		args     []string
		expected res
	}{
		{"--help", "-h", []string{"./testcli", "-h"}, res{"", true}},
		{"--help", "-h", []string{"./testcli", "--help"}, res{"", true}},
		{"--help", "", []string{"./testcli", "--help"}, res{"", true}},
		{"--help", "-h", []string{"./testcli"}, res{"", false}},
		{"--config", "-f", []string{"./testcli", "-f", "/test/config"}, res{"/test/config", true}},
		{"--config", "-f", []string{"./testcli", "--config", "/test/config"}, res{"/test/config", true}},
		{"--config", "-f", []string{"./testcli", "-f"}, res{"", true}},
		{"--config", "-f", []string{"./testcli"}, res{"", false}},
	}

	for _, v := range tests {
		t.Run("with flag " + v.long, func(t *testing.T) {
			os.Args = v.args
			value, ok := RetrieveFlagFromCLI(v.long, v.short)
			assert.Equal(t, v.expected.value, value)
			assert.Equal(t, v.expected.ok, ok)
		})
	}
}
