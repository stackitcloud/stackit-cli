package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type enumSliceFlag struct {
	ignoreCase bool
	options    []string
	value      []string
	valueSet   bool
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &enumFlag{}

// EnumSliceFlag returns a flag which is a slice which values must be one of the given values.
// If ignoreCase is true, values are returned in lower case.
func EnumSliceFlag(ignoreCase bool, defaultValues []string, options ...string) *enumSliceFlag {
	f := &enumSliceFlag{ignoreCase: ignoreCase, options: options}
	err := f.appendToValue(defaultValues)
	if err != nil {
		panic(err)
	}
	return f
}

func (f *enumSliceFlag) appendToValue(values []string) error {
	for _, v := range values {
		v = strings.TrimSpace(v)

		foundValid := false
		for _, o := range f.options {
			if !f.ignoreCase && v == o {
				f.value = append(f.value, v)
				foundValid = true
				break
			} else if f.ignoreCase && strings.EqualFold(v, o) {
				f.value = append(f.value, strings.ToLower(v))
				foundValid = true
				break
			}
		}

		if !foundValid {
			return fmt.Errorf("found value %q, expected one of %q", v, f.options)
		}
	}
	return nil
}

func (f *enumSliceFlag) String() string {
	return "[" + strings.Join(f.value, ",") + "]"
}

func (f *enumSliceFlag) Set(value string) error {
	// If the default value is still set, remove it
	// (Since we're going to append the incoming values to f.value)
	if !f.valueSet {
		f.value = []string{}
		f.valueSet = true
	}

	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	values := strings.Split(value, ",")
	return f.appendToValue(values)
}

func (f *enumSliceFlag) Type() string {
	return "stringSlice"
}
