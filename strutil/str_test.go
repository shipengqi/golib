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

func TestDeDuplicateStr(t *testing.T) {
	tests := []struct {
		title    string
		input    []string
		expected []string
	}{
		{"length should be 1", []string{"111", "111", ""}, []string{"111"}},
		{"length should be 2", []string{"111s", "222", "111s"}, []string{"111s", "222"}},
	}

	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			got := DeDuplicateStr(v.input)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestString2Int(t *testing.T) {
	tests := []struct {
		title    string
		input    []string
		expected []int
		err      bool
	}{
		{"convert string slice to int slice", []string{"111", "121", "222"}, []int{111, 121, 222}, false},
		{"convert error with abc", []string{"111", "abc", "222"}, nil, true},
		{"convert error with empty string", []string{"111", "", "222"}, nil, true},
	}

	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			got, err := String2Int(v.input)
			if v.err {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestContainsOnly(t *testing.T) {
	tests := []struct {
		title    string
		input    string
		input2   string
		expected bool
	}{
		{"should be true", "234234", "0123456789", true},
		{"should be false", "as234234", "0123456789", false},
		{"should be true", "234234%", "0123456789", false},
		{"should be true", "testcase2323-", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789", true},
		{"should be true", "testcase2323", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789", true},
		{"should be true", "testcase", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789", true},
		{"should be false", "testcase+", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789", false},
		{"should be true", "testcase2323+", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789", false},
		{"should be true", "testcase 2323", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789", false},
	}

	for _, v := range tests {
		t.Run(fmt.Sprintf("'%s' contains only '%s' %s", v.input, v.input2, v.title), func(t *testing.T) {
			got := ContainsOnly(v.input, v.input2)
			assert.Equal(t, v.expected, got)
		})
	}
}
