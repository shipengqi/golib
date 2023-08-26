package sysutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBigEndian(t *testing.T) {
	got := IsBigEndian()
	assert.False(t, got)
}

func TestHomeDir(t *testing.T) {
	got := HomeDir()
	assert.NotEmpty(t, got)
}

func TestFQDN(t *testing.T) {
	got, err := FQDN()
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestPwd(t *testing.T) {
	got := Pwd()
	assert.NotEmpty(t, got)
}

func TestIsLinux(t *testing.T) {
	got := IsLinux()
	if isci() {
		assert.True(t, got)
	}
}

func TestIsWindows(t *testing.T) {
	got := IsWindows()
	if isci() {
		assert.False(t, got)
	}
}

func TestUser(t *testing.T) {
	got := User()
	assert.NotEmpty(t, got)
}

func TestTmpDir(t *testing.T) {
	got := TmpDir()
	assert.NotEmpty(t, got)
}

func TestEnvOr(t *testing.T) {
	testEnvKey := "TEST_ENV_KEY"
	testEnvValue := "test"
	unknownEnvKey := "UNKNOWN_ENV_KEY"
	unknownEnvValue := "unknown"
	err := os.Setenv(testEnvKey, testEnvValue)
	assert.NoError(t, err)
	got := EnvOr(testEnvKey, "")
	assert.Equal(t, testEnvValue, got)

	got = EnvOr(unknownEnvKey, unknownEnvValue)
	assert.Equal(t, unknownEnvValue, got)
}

func isci() bool {
	if os.Getenv("CI") == "true" {
		return true
	}
	return false
}
