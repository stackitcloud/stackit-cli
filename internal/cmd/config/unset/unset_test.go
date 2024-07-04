package unset

import (
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
)

func fixtureFlagValues(mods ...func(flagValues map[string]bool)) map[string]bool {
	flagValues := map[string]bool{
		asyncFlag:        true,
		outputFormatFlag: true,
		projectIdFlag:    true,
		verbosityFlag:    true,

		sessionTimeLimitFlag:               true,
		identityProviderCustomEndpointFlag: true,

		argusCustomEndpointFlag:           true,
		authorizationCustomEndpointFlag:   true,
		dnsCustomEndpointFlag:             true,
		loadBalancerCustomEndpointFlag:    true,
		logMeCustomEndpointFlag:           true,
		mariaDBCustomEndpointFlag:         true,
		objectStorageCustomEndpointFlag:   true,
		openSearchCustomEndpointFlag:      true,
		rabbitMQCustomEndpointFlag:        true,
		redisCustomEndpointFlag:           true,
		resourceManagerCustomEndpointFlag: true,
		secretsManagerCustomEndpointFlag:  true,
		serviceAccountCustomEndpointFlag:  true,
		serverBackupCustomEndpointFlag:    true,
		skeCustomEndpointFlag:             true,
		sqlServerFlexCustomEndpointFlag:   true,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		Async:        true,
		OutputFormat: true,
		ProjectId:    true,
		Verbosity:    true,

		SessionTimeLimit:               true,
		IdentityProviderCustomEndpoint: true,

		ArgusCustomEndpoint:           true,
		AuthorizationCustomEndpoint:   true,
		DNSCustomEndpoint:             true,
		LoadBalancerCustomEndpoint:    true,
		LogMeCustomEndpoint:           true,
		MariaDBCustomEndpoint:         true,
		ObjectStorageCustomEndpoint:   true,
		OpenSearchCustomEndpoint:      true,
		RabbitMQCustomEndpoint:        true,
		RedisCustomEndpoint:           true,
		ResourceManagerCustomEndpoint: true,
		SecretsManagerCustomEndpoint:  true,
		ServiceAccountCustomEndpoint:  true,
		ServerBackupCustomEndpoint:    true,
		SKECustomEndpoint:             true,
		SQLServerFlexCustomEndpoint:   true,
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
				model.Verbosity = false

				model.SessionTimeLimit = false
				model.IdentityProviderCustomEndpoint = false

				model.ArgusCustomEndpoint = false
				model.AuthorizationCustomEndpoint = false
				model.DNSCustomEndpoint = false
				model.LoadBalancerCustomEndpoint = false
				model.LogMeCustomEndpoint = false
				model.MariaDBCustomEndpoint = false
				model.ObjectStorageCustomEndpoint = false
				model.OpenSearchCustomEndpoint = false
				model.RabbitMQCustomEndpoint = false
				model.RedisCustomEndpoint = false
				model.ResourceManagerCustomEndpoint = false
				model.SecretsManagerCustomEndpoint = false
				model.ServiceAccountCustomEndpoint = false
				model.ServerBackupCustomEndpoint = false
				model.SKECustomEndpoint = false
				model.SQLServerFlexCustomEndpoint = false
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
			description: "identity provider custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[identityProviderCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IdentityProviderCustomEndpoint = false
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
		{
			description: "serverbackup custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[serverBackupCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ServerBackupCustomEndpoint = false
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)

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
