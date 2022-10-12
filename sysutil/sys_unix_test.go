package sysutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPID(t *testing.T) {
	assert.True(t, PID() > 0)
}

func TestGetProcessByPid(t *testing.T) {
	got := GetProcessByPid(0)
	assert.Empty(t, got)
}
