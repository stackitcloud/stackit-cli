package flags

import (
	"strings"

	"github.com/spf13/pflag"
)

type cidrSliceFlag struct {
	value []string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &cidrFlag{}

// CIDRSliceFlag returns a flag which must be a valid CIDR slice.
func CIDRSliceFlag() *cidrSliceFlag {
	return &cidrSliceFlag{}
}

func (f *cidrSliceFlag) String() string {
	return "[" + strings.Join(f.value, ",") + "]"
}

func (f *cidrSliceFlag) Set(value string) error {
	if value == "" {
		f.value = []string{}
		return nil
	}

	cidrs := strings.Split(value, ",")

	for i, cidr := range cidrs {
		cidrs[i] = strings.TrimSpace(cidr)

		err := validateCIDR(cidrs[i])
		if err != nil {
			return err
		}
	}

	f.value = append(f.value, cidrs...)
	return nil
}

func (f *cidrSliceFlag) Type() string {
	return "stringSlice"
}
