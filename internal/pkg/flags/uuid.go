package flags

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

type uuidFlag struct {
	value string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &uuidFlag{}

// UUIDFlag returns a flag which must be a valid UUID.
func UUIDFlag() *uuidFlag {
	return &uuidFlag{}
}

func (f *uuidFlag) String() string {
	return f.value
}

func (f *uuidFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	if _, err := uuid.Parse(value); err != nil {
		return fmt.Errorf("parse as UUID: %w", err)
	}
	f.value = value
	return nil
}

func (f *uuidFlag) Type() string {
	return "string"
}
