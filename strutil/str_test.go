package strutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		title    string
		input    string
		input2   string
		expected bool
	}{
		{"should be true", "test lower case", "LoWer", true},
		{"should be false", "test lower case", "LLL", false},
		{"should be true", "TEST UPPER CASE", "upper", true},
	}

	for _, v := range tests {
		t.Run(fmt.Sprintf("'%s' contains '%s' %s", v.input, v.input2, v.title), func(t *testing.T) {
			got := ContainsIgnoreCase(v.input, v.input2)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestEqualsIgnoreCase(t *testing.T) {
	tests := []struct {
		title    string
		input    string
		input2   string
		expected bool
	}{
		{"should be true", "lower", "LoWer", true},
		{"should be true", "lower", "LOWER", true},
		{"should be true", "UPPER", "upper", true},
		{"should be false", "UPPER", "u", false},
	}

	for _, v := range tests {
		t.Run(fmt.Sprintf("'%s' equals '%s' %s", v.input, v.input2, v.title), func(t *testing.T) {
			got := EqualsIgnoreCase(v.input, v.input2)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		title    string
		input    string
		expected bool
	}{
		{"should be empty", "", true},
		{"blank string should be empty", "   ", true},
		{"line break should be empty", "  \n ", true},
	}

	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			got := IsEmpty(v.input)
			assert.Equal(t, v.expected, got)
		})
	}
}
