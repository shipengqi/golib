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
		os.Args = v.args
		value, ok := RetrieveFlagFromCLI(v.long, v.short)
		assert.Equal(t, v.expected.value, value)
		assert.Equal(t, v.expected.ok, ok)
	}
}

func TestIsHelpCmd(t *testing.T) {
	tests := []struct {
		short    string
		args     []string
		expected bool
	}{
		{"", []string{"./testcli", "help"}, true},
		{"-h", []string{"./testcli", "help"}, true},
		{"-h", []string{"./testcli", "-h"}, true},
		{"-h", []string{"./testcli", "--help"}, true},
		{"", []string{"./testcli", "--help"}, true},
		{"", []string{"./testcli", "-h"}, false},
		{"-h", []string{"./testcli"}, false},
		{"-h", []string{"./testcli", "-f", "/test/config"}, false},
		{"-h", []string{"./testcli", "--config", "/test/config"}, false},
		{"-h", []string{"./testcli", "-f"}, false},
	}

	for _, v := range tests {
		os.Args = v.args
		ok := IsHelpCmd(v.short)
		assert.Equal(t, v.expected, ok)
	}
}

func TestIsVersionCmd(t *testing.T) {
	tests := []struct {
		short    string
		args     []string
		expected bool
	}{
		{"", []string{"./testcli", "version"}, true},
		{"-v", []string{"./testcli", "version"}, true},
		{"-v", []string{"./testcli", "-v"}, true},
		{"-v", []string{"./testcli", "--version"}, true},
		{"", []string{"./testcli", "--version"}, true},
		{"", []string{"./testcli", "-v"}, false},
		{"-v", []string{"./testcli"}, false},
		{"-v", []string{"./testcli", "-f", "/test/config"}, false},
		{"-v", []string{"./testcli", "--config", "/test/config"}, false},
		{"-v", []string{"./testcli", "-f"}, false},
	}

	for _, v := range tests {
		os.Args = v.args
		ok := IsVersionCmd(v.short)
		assert.Equal(t, v.expected, ok)
	}
}
