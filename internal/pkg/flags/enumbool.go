package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type enumBoolFlag struct {
	value string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &enumBoolFlag{}

// enumBoolFlag returns a flag which must be either "true" or "false".
// This is different than an usual bool flag, which doesn't take arguments and is either set or unset.
//
// It's almost identical to EnumFlag(true, "true", "false"), but will return a bool value instead of a string value.
func EnumBoolFlag() *enumBoolFlag {
	return &enumBoolFlag{}
}

func (f *enumBoolFlag) String() string {
	return f.value
}

func (f *enumBoolFlag) Set(value string) error {
	valueLower := strings.ToLower(value)
	if valueLower != "true" && valueLower != "false" {
		return fmt.Errorf("expected one of %q", []string{"true", "false"})
	}
	f.value = valueLower
	return nil
}

func (f *enumBoolFlag) Type() string {
	return "bool"
}
