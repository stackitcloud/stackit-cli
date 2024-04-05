package unset

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func fixtureFlagValues(mods ...func(flagValues map[string]bool)) map[string]bool {
	flagValues := map[string]bool{
		asyncFlag:            true,
		outputFormatFlag:     true,
		projectIdFlag:        true,
		sessionTimeLimitFlag: true,
		verbosityFlag:        true,

		argusCustomEndpointFlag:           true,
		authorizationCustomEndpointFlag:   true,
		dnsCustomEndpointFlag:             true,
		logMeCustomEndpointFlag:           true,
		mariaDBCustomEndpointFlag:         true,
		objectStorageCustomEndpointFlag:   true,
		openSearchCustomEndpointFlag:      true,
		rabbitMQCustomEndpointFlag:        true,
		redisCustomEndpointFlag:           true,
		resourceManagerCustomEndpointFlag: true,
		secretsManagerCustomEndpointFlag:  true,
		serviceAccountCustomEndpointFlag:  true,
		skeCustomEndpointFlag:             true,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		Async:            true,
		OutputFormat:     true,
		ProjectId:        true,
		SessionTimeLimit: true,
		Verbosity:        true,

		ArgusCustomEndpoint:           true,
		AuthorizationCustomEndpoint:   true,
		DNSCustomEndpoint:             true,
		LogMeCustomEndpoint:           true,
		MariaDBCustomEndpoint:         true,
		ObjectStorageCustomEndpoint:   true,
		OpenSearchCustomEndpoint:      true,
		RabbitMQCustomEndpoint:        true,
		RedisCustomEndpoint:           true,
		ResourceManagerCustomEndpoint: true,
		SecretsManagerCustomEndpoint:  true,
		ServiceAccountCustomEndpoint:  true,
		SKECustomEndpoint:             true,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]bool
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
			flagValues:  map[string]bool{},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Async = false
				model.OutputFormat = false
				model.ProjectId = false
				model.SessionTimeLimit = false
				model.Verbosity = false

				model.ArgusCustomEndpoint = false
				model.AuthorizationCustomEndpoint = false
				model.DNSCustomEndpoint = false
				model.LogMeCustomEndpoint = false
				model.MariaDBCustomEndpoint = false
				model.ObjectStorageCustomEndpoint = false
				model.OpenSearchCustomEndpoint = false
				model.RabbitMQCustomEndpoint = false
				model.RedisCustomEndpoint = false
				model.ResourceManagerCustomEndpoint = false
				model.SecretsManagerCustomEndpoint = false
				model.ServiceAccountCustomEndpoint = false
				model.SKECustomEndpoint = false
			}),
		},
		{
			description: "project id empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[projectIdFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProjectId = false
			}),
		},
		{
			description: "output format empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[outputFormatFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.OutputFormat = false
			}),
		},
		{
			description: "argus custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[argusCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ArgusCustomEndpoint = false
			}),
		},
		{
			description: "dns custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[dnsCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DNSCustomEndpoint = false
			}),
		},
		{
			description: "secrets manager custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[secretsManagerCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SecretsManagerCustomEndpoint = false
			}),
		},
		{
			description: "service account custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[serviceAccountCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ServiceAccountCustomEndpoint = false
			}),
		},
		{
			description: "ske custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[skeCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SKECustomEndpoint = false
			}),
		},
		{
			description: "resource manager custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[resourceManagerCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ResourceManagerCustomEndpoint = false
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd(nil)

			for flag, value := range tt.flagValues {
				stringBool := fmt.Sprintf("%v", value)
				err := cmd.Flags().Set(flag, stringBool)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, stringBool, err)
				}
			}

			err := cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model := parseInput(cmd)

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
