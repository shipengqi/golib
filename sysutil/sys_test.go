package sysutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBigEndian(t *testing.T) {
	got := IsBigEndian()
	assert.False(t, got)
}
