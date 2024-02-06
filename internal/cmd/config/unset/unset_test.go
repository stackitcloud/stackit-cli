package unset

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func fixtureFlagValues(mods ...func(flagValues map[string]bool)) map[string]bool {
	flagValues := map[string]bool{
		projectIdFlag:    true,
		outputFormatFlag: true,

		dnsCustomEndpointFlag:             true,
		openSearchCustomEndpointFlag:      true,
		rabbitMQCustomEndpointFlag:        true,
		redisCustomEndpointFlag:           true,
		resourceManagerCustomEndpointFlag: true,
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
		ProjectId:                     true,
		OutputFormat:                  true,
		DNSCustomEndpoint:             true,
		ServiceAccountCustomEndpoint:  true,
		SKECustomEndpoint:             true,
		ResourceManagerCustomEndpoint: true,
		OpenSearchCustomEndpoint:      true,
		RedisCustomEndpoint:           true,
		RabbitMQCustomEndpoint:        true,
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
				model.ProjectId = false
				model.OutputFormat = false
				model.DNSCustomEndpoint = false
				model.ServiceAccountCustomEndpoint = false
				model.SKECustomEndpoint = false
				model.ResourceManagerCustomEndpoint = false
				model.OpenSearchCustomEndpoint = false
				model.RedisCustomEndpoint = false
				model.RabbitMQCustomEndpoint = false
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
			cmd := NewCmd()

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
