package utils

import "fmt"

// PtrString creates a string representation of a passed object pointer or returns
// an empty string, if the passed object is _nil_.
func PtrString[T any](t *T) string {
	if t != nil {
		return fmt.Sprintf("%v", *t)
	}
	return ""
}
