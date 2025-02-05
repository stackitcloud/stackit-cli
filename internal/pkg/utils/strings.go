package utils

import (
	"strings"
)

// JoinStringKeys concatenates the string keys of a map, each separatore by the
// [sep] string.
func JoinStringKeys(m map[string]any, sep string) string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return strings.Join(keys, sep)
}

// JoinStringKeysPtr concatenates the string keys of a map pointer, each separatore by the
// [sep] string.
func JoinStringKeysPtr(m map[string]any, sep string) string {
	if m == nil {
		return ""
	}
	return JoinStringKeys(m, sep)
}

// JoinStringPtr concatenates the strings of a string slice pointer, each separatore by the
// [sep] string.
func JoinStringPtr(vals *[]string, sep string) string {
	if vals == nil || len(*vals) == 0 {
		return ""
	}
	return strings.Join(*vals, sep)
}
