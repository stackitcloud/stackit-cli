package set

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		description string
		flagValues  map[string]string
		isValid     bool
	}{
		{
			description: "valid session time limit 1",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "1h",
			},
			isValid: true,
		},
		{
			description: "valid session time limit 2",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "5h30m40s",
			},
			isValid: true,
		},
		{
			description: "valid session time limit 3",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "1h2m3s4ms5us6ns",
			},
			isValid: true,
		},
		{
			description: "valid session time limit 4",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "1d",
			},
			isValid: true,
		},
		{
			description: "invalid session time limit 1",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "foo",
			},
			isValid: false,
		},
		{
			description: "invalid session time limit 2",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "",
			},
			isValid: false,
		},
		{
			description: "invalid session time limit 3",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "1",
			},
			isValid: false,
		},
		{
			description: "invalid session time limit 4",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "h",
			},
			isValid: false,
		},
		{
			description: "invalid session time limit 5",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "0h",
			},
			isValid: false,
		},
		{
			description: "invalid session time limit 6",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "-1h",
			},
			isValid: false,
		},
		{
			description: "invalid session time limit 7",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "25h",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}

			configureFlags(cmd)

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err := cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}
			_, err = parseFlags(cmd)

			if err != nil && tt.isValid {
				t.Fatalf("error parsing flags: %v", err)
			}
			if err == nil && !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
		})
	}
}
