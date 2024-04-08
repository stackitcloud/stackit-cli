package args

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

func TestNoArgs(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		isValid     bool
	}{
		{
			description: "valid",
			args:        nil,
			isValid:     true,
		},
		{
			description: "invalid",
			args:        []string{"unknown"},
			isValid:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:   "test",
				Short: "Test command",
			}

			err := NoArgs(cmd, tt.args)

			if tt.isValid && err != nil {
				t.Fatalf("should not have failed: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("should have failed")
			}
		})
	}
}

func TestSingleArg(t *testing.T) {
	tests := []struct {
		description  string
		args         []string
		validateFunc func(value string) error
		isValid      bool
	}{
		{
			description: "valid",
			args:        []string{"arg"},
			validateFunc: func(value string) error {
				return nil
			},
			isValid: true,
		},
		{
			description: "no_arg",
			args:        []string{},
			isValid:     false,
		},
		{
			description: "more_than_one_arg",
			args:        []string{"arg", "arg2"},
			isValid:     false,
		},
		{
			description: "empty_arg",
			args:        []string{""},
			isValid:     false,
		},
		{
			description: "invalid_arg",
			args:        []string{"arg"},
			validateFunc: func(value string) error {
				return fmt.Errorf("error")
			},
			isValid: false,
		},
		{
			description:  "nil validation function",
			args:         []string{"arg"},
			validateFunc: nil,
			isValid:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:   "test",
				Short: "Test command",
			}

			argFunction := SingleArg("test", tt.validateFunc)
			err := argFunction(cmd, tt.args)

			if tt.isValid && err != nil {
				t.Fatalf("should not have failed: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("should have failed")
			}
		})
	}
}

func TestSingleOptionalArg(t *testing.T) {
	tests := []struct {
		description  string
		args         []string
		validateFunc func(value string) error
		isValid      bool
	}{
		{
			description: "valid",
			args:        []string{"arg"},
			validateFunc: func(value string) error {
				return nil
			},
			isValid: true,
		},
		{
			description: "no_arg",
			args:        []string{},
			isValid:     true,
		},
		{
			description: "more_than_one_arg",
			args:        []string{"arg", "arg2"},
			isValid:     false,
		},
		{
			description: "empty_arg",
			args:        []string{""},
			isValid:     true,
		},
		{
			description: "invalid_arg",
			args:        []string{"arg"},
			validateFunc: func(value string) error {
				return fmt.Errorf("error")
			},
			isValid: false,
		},
		{
			description:  "nil validation function",
			args:         []string{"arg"},
			validateFunc: nil,
			isValid:      true,
		},
		{
			description:  "nil validation function, no args",
			args:         []string{},
			validateFunc: nil,
			isValid:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:   "test",
				Short: "Test command",
			}

			argFunction := SingleOptionalArg("test", tt.validateFunc)
			err := argFunction(cmd, tt.args)

			if tt.isValid && err != nil {
				t.Fatalf("should not have failed: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("should have failed")
			}
		})
	}
}
