package activateserviceaccount

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/zalando/go-keyring"

	"github.com/google/go-cmp/cmp"
)

var testTokenCustomEndpoint = "token_url"
var testJwksCustomEndpoint = "jwks_url"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		serviceAccountTokenFlag:   "token",
		serviceAccountKeyPathFlag: "sa_key",
		privateKeyPathFlag:        "private_key",
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
		description         string
		flagValues          map[string]string
		tokenCustomEndpoint string
		jwksCustomEndpoint  string
		isValid             bool
		expectedModel       *inputModel
	}{
		{
			description:         "base",
			flagValues:          fixtureFlagValues(),
			tokenCustomEndpoint: testTokenCustomEndpoint,
			jwksCustomEndpoint:  testJwksCustomEndpoint,
			isValid:             true,
			expectedModel:       fixtureInputModel(),
		},
		{
			description:         "no values",
			flagValues:          map[string]string{},
			tokenCustomEndpoint: "",
			jwksCustomEndpoint:  "",
			isValid:             true,
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
			},
			tokenCustomEndpoint: "",
			jwksCustomEndpoint:  "",
			isValid:             true,
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
			p := print.NewPrinter()
			cmd := NewCmd(p)
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

			viper.Reset()
			viper.Set(config.TokenCustomEndpointKey, tt.tokenCustomEndpoint)
			viper.Set(config.JwksCustomEndpointKey, tt.jwksCustomEndpoint)

			model := parseInput(p, cmd)

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

func TestStoreFlags(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		isValid     bool
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			isValid:     true,
		},
		{
			description: "no values",
			model: &inputModel{
				ServiceAccountToken:   "",
				ServiceAccountKeyPath: "",
				PrivateKeyPath:        "",
				TokenCustomEndpoint:   "",
				JwksCustomEndpoint:    "",
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Initialize an empty keyring
			keyring.MockInit()

			err := storeFlags(tt.model)
			if !tt.isValid {
				if err == nil {
					t.Fatalf("did not fail on invalid input")
				}
				return
			}
			if err != nil {
				t.Fatalf("store flags: %v", err)
			}

			value, err := auth.GetAuthField(auth.TOKEN_CUSTOM_ENDPOINT)
			if err != nil {
				t.Errorf("Failed to get value of auth field: %v", err)
			}
			if value != tt.model.TokenCustomEndpoint {
				t.Errorf("Value of \"%s\" does not match: expected \"%s\", got \"%s\"", auth.TOKEN_CUSTOM_ENDPOINT, tt.model.TokenCustomEndpoint, value)
			}

			value, err = auth.GetAuthField(auth.JWKS_CUSTOM_ENDPOINT)
			if err != nil {
				t.Errorf("Failed to get value of auth field: %v", err)
			}
			if value != tt.model.JwksCustomEndpoint {
				t.Errorf("Value of \"%s\" does not match: expected \"%s\", got \"%s\"", auth.JWKS_CUSTOM_ENDPOINT, tt.model.TokenCustomEndpoint, value)
			}
		})
	}
}
