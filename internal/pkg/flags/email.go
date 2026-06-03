package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type emailFlag struct {
	value string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &emailFlag{}

// EmailFlag returns a flag which must be a valid Email.
func EmailFlag() *emailFlag {
	return &emailFlag{}
}

func (f *emailFlag) String() string {
	return f.value
}

func (f *emailFlag) Set(value string) error {
	isEmail := value != "" && strings.Contains(value, "@")
	if !isEmail {
		return fmt.Errorf("invalid email address: %s", value)
	}
	f.value = value
	return nil
}

func (f *emailFlag) Type() string {
	return "string"
}
