package flags

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestEmailFlag(t *testing.T) {
	tests := []struct {
		description string
		value       string
		isValid     bool
	}{
		{
			description: "valid",
			value:       "test@test",
			isValid:     true,
		},
		{
			description: "empty",
			value:       "",
			isValid:     false,
		},
		{
			description: "invalid",
			value:       "invalid-email",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			flag := EmailFlag()
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(_ *cobra.Command, _ []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			err := cmd.Flags().Set("test-flag", tt.value)

			if !tt.isValid && err == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if err != nil {
				t.Fatalf("failed on valid input: %v", err)
			}
			value := FlagToStringValue(nil, cmd, "test-flag")
			if value != tt.value {
				t.Fatalf("flag did not return set value")
			}
		})
	}
}
