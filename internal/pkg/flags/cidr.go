package flags

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type cidrFlag struct {
	value string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &cidrFlag{}

// CIDRFlag returns a flag which must be a valid CIDR.
func CIDRFlag() *cidrFlag {
	return &cidrFlag{}
}

func (f *cidrFlag) String() string {
	return f.value
}

func (f *cidrFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	err := validateCIDR(value)
	if err != nil {
		return err
	}
	f.value = value
	return nil
}

func (f *cidrFlag) Type() string {
	return "string"
}

func validateCIDR(value string) error {
	_, _, err := net.ParseCIDR(value)
	if err != nil {
		return fmt.Errorf("parse %s as CIDR: %w", value, err)
	}
	return nil
}
