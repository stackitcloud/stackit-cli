package set

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description: "valid session time limit 1",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "1h",
			},
			isValid: true,
			expectedModel: &inputModel{
				SessionTimeLimit: utils.Ptr("1h"),
			},
		},
		{
			description: "valid session time limit 2",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "5h30m40s",
			},
			isValid: true,
			expectedModel: &inputModel{
				SessionTimeLimit: utils.Ptr("5h30m40s"),
			},
		},
		{
			description: "valid session time limit 3",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "1h2m3s4ms5us6ns",
			},
			isValid: true,
			expectedModel: &inputModel{
				SessionTimeLimit: utils.Ptr("1h2m3s4ms5us6ns"),
			},
		},
		{
			description: "valid session time limit 4",
			flagValues: map[string]string{
				sessionTimeLimitFlag: "1d",
			},
			isValid: true,
			expectedModel: &inputModel{
				SessionTimeLimit: utils.Ptr("24h"),
			},
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
		{
			description: "project ID set",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: uuid.NewString(),
			},
			isValid: true,
			expectedModel: &inputModel{
				ProjectIdSet: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd(nil)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(cmd, nil)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
