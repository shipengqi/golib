package strutil

import (
	"strconv"
	"strings"
)

// EqualsIgnoreCase check if str1 is equal to str2 ignoring case sensitivity
func EqualsIgnoreCase(str1, str2 string) bool {
	// SA6005 - Inefficient string comparison with strings.ToLower or strings.ToUpper
	// if strings.ToLower(s1) == strings.ToLower(s2) { ... }
	return strings.EqualFold(str1, str2)
}

// ContainsIgnoreCase check if str1 contains str2 ignoring case sensitivity
func ContainsIgnoreCase(str1, str2 string) bool {
	return strings.Contains(strings.ToLower(str1), strings.ToLower(str2))
}

// IsEmpty check if value is empty string or blank string.
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// DeDuplicateStr De-Duplicate the given strings.
func DeDuplicateStr(s []string) []string {
	encountered := map[string]struct{}{}
	ret := make([]string, 0)
	for i := range s {
		if len(s[i]) == 0 {
			continue
		}
		if _, contained := encountered[s[i]]; contained {
			continue
		}
		encountered[s[i]] = struct{}{}
		ret = append(ret, s[i])
	}
	return ret
}

// String2Int convert string slice to int slice.
func String2Int(arr []string) ([]int, error) {
	var err error
	ints := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		ints[i], err = strconv.Atoi(arr[i])
		if err != nil {
			return nil, err
		}
	}
	return ints, nil
}
