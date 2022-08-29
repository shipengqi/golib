package strutil

import "strings"

// EqualsIgnoreCase check if str1 is equal to str2 ignoring case sensitivity
func EqualsIgnoreCase(str1, str2 string) bool {
	return strings.ToLower(str1) == strings.ToLower(str2)
}

// ContainsIgnoreCase check if str1 contains str2 ignoring case sensitivity
func ContainsIgnoreCase(str1, str2 string) bool {
	return strings.Contains(strings.ToLower(str1), strings.ToLower(str2))
}

// IsEmpty check if value is empty string or blank string.
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
