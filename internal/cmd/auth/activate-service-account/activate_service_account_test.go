package activateserviceaccount

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

var testTokenCustomEndpoint = "token_url"

func boolPtr(v bool) *bool {
	return &v
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		serviceAccountTokenFlag:   "token",
		serviceAccountKeyPathFlag: "sa_key",
		privateKeyPathFlag:        "private_key",
		useOIDCFlag:               "true",
		onlyPrintAccessTokenFlag:  "true",
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
		UseOIDC:               boolPtr(true),
		OnlyPrintAccessToken:  true,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description         string
		argValues           []string
		flagValues          map[string]string
		tokenCustomEndpoint string
		isValid             bool
		expectedModel       *inputModel
	}{
		{
			description:         "base",
			flagValues:          fixtureFlagValues(),
			tokenCustomEndpoint: testTokenCustomEndpoint,
			isValid:             true,
			expectedModel:       fixtureInputModel(),
		},
		{
			description:         "no values",
			flagValues:          map[string]string{},
			tokenCustomEndpoint: "",
			isValid:             true,
			expectedModel: &inputModel{
				ServiceAccountToken:   "",
				ServiceAccountKeyPath: "",
				PrivateKeyPath:        "",
				UseOIDC:               nil,
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
			isValid:             true,
			expectedModel: &inputModel{
				ServiceAccountToken:   "",
				ServiceAccountKeyPath: "",
				PrivateKeyPath:        "",
				UseOIDC:               nil,
			},
		},
		{
			description: "use_oidc_true",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[useOIDCFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.UseOIDC = boolPtr(true)
			}),
		},
		{
			description: "use_oidc_false",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[useOIDCFlag] = "false"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.UseOIDC = boolPtr(false)
			}),
		},
		{
			description: "invalid_flag",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["test_flag"] = "test"
			}),
			isValid: false,
		},
		{
			description: "default value OnlyPrintAccessToken",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, "only-print-access-token")
				},
			),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.OnlyPrintAccessToken = false
			}),
		},
		{
			description: "default value UseOIDC",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, useOIDCFlag)
				},
			),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.UseOIDC = nil
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestStoreCustomEndpointFlags(t *testing.T) {
	tests := []struct {
		description         string
		model               *inputModel
		tokenCustomEndpoint string
		isValid             bool
	}{
		{
			description:         "base",
			model:               fixtureInputModel(),
			tokenCustomEndpoint: testTokenCustomEndpoint,
			isValid:             true,
		},
		{
			description: "no values",
			model: &inputModel{
				ServiceAccountToken:   "",
				ServiceAccountKeyPath: "",
				PrivateKeyPath:        "",
			},
			tokenCustomEndpoint: "",
			isValid:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Initialize an empty keyring
			keyring.MockInit()

			viper.Reset()
			viper.Set(config.TokenCustomEndpointKey, tt.tokenCustomEndpoint)

			err := storeCustomEndpoint(tt.tokenCustomEndpoint)
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
			if value != tt.tokenCustomEndpoint {
				t.Errorf("Value of \"%s\" does not match: expected \"%s\", got \"%s\"", auth.TOKEN_CUSTOM_ENDPOINT, tt.tokenCustomEndpoint, value)
			}
		})
	}
}
