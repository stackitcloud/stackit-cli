package activateserviceaccount

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"

	"github.com/google/go-cmp/cmp"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		serviceAccountTokenFlag:   "token",
		serviceAccountKeyPathFlag: "sa_key",
		privateKeyPathFlag:        "private_key",
		tokenCustomEndpointFlag:   "token_url",
		jwksCustomEndpointFlag:    "jwks_url",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		ServiceAccountToken:   "token",
		ServiceAccountKeyPath: "sa_key",
		PrivateKeyPath:        "private_key",
		TokenCustomEndpoint:   "token_url",
		JwksCustomEndpoint:    "jwks_url",
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: &inputModel{
				ServiceAccountToken:   "",
				ServiceAccountKeyPath: "",
				PrivateKeyPath:        "",
				TokenCustomEndpoint:   "",
				JwksCustomEndpoint:    "",
			},
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				serviceAccountTokenFlag:   "",
				serviceAccountKeyPathFlag: "",
				privateKeyPathFlag:        "",
				tokenCustomEndpointFlag:   "",
				jwksCustomEndpointFlag:    "",
			},
			isValid: true,
			expectedModel: &inputModel{
				ServiceAccountToken:   "",
				ServiceAccountKeyPath: "",
				PrivateKeyPath:        "",
				TokenCustomEndpoint:   "",
				JwksCustomEndpoint:    "",
			},
		},
		{
			description: "invalid_flag",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["test_flag"] = "test"
			}),
			isValid: false,
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

			model := parseInput(cmd, nil)

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
