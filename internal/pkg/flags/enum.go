package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type enumFlag struct {
	ignoreCase bool
	options    []string
	value      string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &enumFlag{}

// EnumFlag returns a flag which must be one of the given values.
// If ignoreCase is true, flag value is returned in lower case.
func EnumFlag(ignoreCase bool, options ...string) *enumFlag {
	return &enumFlag{ignoreCase: ignoreCase, options: options}
}

func (f *enumFlag) String() string {
	return f.value
}

func (f *enumFlag) Set(value string) error {
	for _, o := range f.options {
		if !f.ignoreCase && value == o {
			f.value = value
			return nil
		}
		if f.ignoreCase && strings.EqualFold(value, o) {
			f.value = strings.ToLower(value)
			return nil
		}
	}

	return fmt.Errorf("expected one of %q", f.options)
}

func (f *enumFlag) Type() string {
	return "string"
}
