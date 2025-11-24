package unset

import (
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
)

func fixtureFlagValues(mods ...func(flagValues map[string]bool)) map[string]bool {
	flagValues := map[string]bool{
		asyncFlag:        true,
		outputFormatFlag: true,
		projectIdFlag:    true,
		verbosityFlag:    true,

		sessionTimeLimitFlag:                             true,
		identityProviderCustomWellKnownConfigurationFlag: true,
		identityProviderCustomClientIdFlag:               true,
		allowedUrlDomainFlag:                             true,

		authorizationCustomEndpointFlag:   true,
		dnsCustomEndpointFlag:             true,
		loadBalancerCustomEndpointFlag:    true,
		logMeCustomEndpointFlag:           true,
		mariaDBCustomEndpointFlag:         true,
		objectStorageCustomEndpointFlag:   true,
		observabilityCustomEndpointFlag:   true,
		openSearchCustomEndpointFlag:      true,
		rabbitMQCustomEndpointFlag:        true,
		redisCustomEndpointFlag:           true,
		resourceManagerCustomEndpointFlag: true,
		secretsManagerCustomEndpointFlag:  true,
		kmsCustomEndpointFlag:             true,
		serviceAccountCustomEndpointFlag:  true,
		serverBackupCustomEndpointFlag:    true,
		serverOsUpdateCustomEndpointFlag:  true,
		runCommandCustomEndpointFlag:      true,
		skeCustomEndpointFlag:             true,
		sqlServerFlexCustomEndpointFlag:   true,
		iaasCustomEndpointFlag:            true,
		tokenCustomEndpointFlag:           true,
		intakeCustomEndpointFlag:          true,
		cdnCustomEndpointFlag:             true,
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
		IdentityProviderCustomClientID: true,
		AllowedUrlDomain:               true,

		AuthorizationCustomEndpoint:   true,
		DNSCustomEndpoint:             true,
		LoadBalancerCustomEndpoint:    true,
		LogMeCustomEndpoint:           true,
		MariaDBCustomEndpoint:         true,
		ObjectStorageCustomEndpoint:   true,
		ObservabilityCustomEndpoint:   true,
		OpenSearchCustomEndpoint:      true,
		RabbitMQCustomEndpoint:        true,
		RedisCustomEndpoint:           true,
		ResourceManagerCustomEndpoint: true,
		SecretsManagerCustomEndpoint:  true,
		KMSCustomEndpoint:             true,
		ServiceAccountCustomEndpoint:  true,
		ServerBackupCustomEndpoint:    true,
		ServerOsUpdateCustomEndpoint:  true,
		RunCommandCustomEndpoint:      true,
		SKECustomEndpoint:             true,
		SQLServerFlexCustomEndpoint:   true,
		IaaSCustomEndpoint:            true,
		TokenCustomEndpoint:           true,
		IntakeCustomEndpoint:          true,
		CDNCustomEndpoint:             true,
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
				model.IdentityProviderCustomClientID = false
				model.AllowedUrlDomain = false

				model.AuthorizationCustomEndpoint = false
				model.DNSCustomEndpoint = false
				model.LoadBalancerCustomEndpoint = false
				model.LogMeCustomEndpoint = false
				model.MariaDBCustomEndpoint = false
				model.ObjectStorageCustomEndpoint = false
				model.ObservabilityCustomEndpoint = false
				model.OpenSearchCustomEndpoint = false
				model.RabbitMQCustomEndpoint = false
				model.RedisCustomEndpoint = false
				model.ResourceManagerCustomEndpoint = false
				model.SecretsManagerCustomEndpoint = false
				model.KMSCustomEndpoint = false
				model.ServiceAccountCustomEndpoint = false
				model.ServerBackupCustomEndpoint = false
				model.ServerOsUpdateCustomEndpoint = false
				model.RunCommandCustomEndpoint = false
				model.SKECustomEndpoint = false
				model.SQLServerFlexCustomEndpoint = false
				model.IaaSCustomEndpoint = false
				model.TokenCustomEndpoint = false
				model.IntakeCustomEndpoint = false
				model.CDNCustomEndpoint = false
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
				flagValues[identityProviderCustomWellKnownConfigurationFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IdentityProviderCustomEndpoint = false
			}),
		},
		{
			description: "identity provider custom client id empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[identityProviderCustomClientIdFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IdentityProviderCustomClientID = false
			}),
		},
		{
			description: "allowed url domain empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[allowedUrlDomainFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.AllowedUrlDomain = false
			}),
		},
		{
			description: "observability custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[observabilityCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ObservabilityCustomEndpoint = false
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
			description: "kms custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[kmsCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.KMSCustomEndpoint = false
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
		{
			description: "serverosupdate custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[serverOsUpdateCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ServerOsUpdateCustomEndpoint = false
			}),
		},
		{
			description: "runcommand custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[runCommandCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.RunCommandCustomEndpoint = false
			}),
		},
		{
			description: "token custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[tokenCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.TokenCustomEndpoint = false
			}),
		},
		{
			description: "cdn custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[cdnCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.CDNCustomEndpoint = false
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})

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
