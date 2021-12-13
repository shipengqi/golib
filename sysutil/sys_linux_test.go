package sysutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPID(t *testing.T) {
	assert.True(t, PID() > 0)
}
