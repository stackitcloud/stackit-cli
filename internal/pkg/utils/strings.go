package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
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

// JoinStringMap concatenates the key-value pairs of a string map, key and value separated by kvSep, key value pairs separated by sep.
func JoinStringMap(m map[string]string, kvSep, sep string) string {
	if m == nil {
		return ""
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s%s%s", k, kvSep, v))
	}
	return strings.Join(parts, sep)
}

// JoinStringPtr concatenates the strings of a string slice pointer, each separatore by the
// [sep] string.
func JoinStringPtr(vals *[]string, sep string) string {
	if vals == nil || len(*vals) == 0 {
		return ""
	}
	return strings.Join(*vals, sep)
}

// Truncate trims the passed string (if it is not nil). If the input string is
// longer than the given length, it is truncated to _maxLen_ and a ellipsis (…)
// is attached. Therefore the resulting string has at most length _maxLen-1_
func Truncate(s *string, maxLen int) string {
	if s == nil {
		return ""
	}

	if utf8.RuneCountInString(*s) > maxLen {
		var builder strings.Builder
		for i, r := range *s {
			if i >= maxLen {
				break
			}
			builder.WriteRune(r)
		}
		builder.WriteRune('…')
		return builder.String()
	}
	return *s
}
