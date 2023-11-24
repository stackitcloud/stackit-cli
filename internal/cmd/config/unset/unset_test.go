package unset

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func fixtureFlagValues(mods ...func(flagValues map[string]bool)) map[string]bool {
	flagValues := map[string]bool{
		projectIdFlag:                true,
		outputFormatFlag:             true,
		dnsCustomEndpointFlag:        true,
		postgreSQLCustomEndpointFlag: true,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagModel(mods ...func(model *flagModel)) *flagModel {
	model := &flagModel{
		ProjectId:                true,
		OutputFormat:             true,
		DNSCustomEndpoint:        true,
		PostgreSQLCustomEndpoint: true,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]bool
		isValid       bool
		expectedModel *flagModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureFlagModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]bool{},
			isValid:     true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.ProjectId = false
				model.OutputFormat = false
				model.DNSCustomEndpoint = false
				model.PostgreSQLCustomEndpoint = false
			}),
		},
		{
			description: "project id empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[projectIdFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.ProjectId = false
			}),
		},
		{
			description: "output format empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[outputFormatFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.OutputFormat = false
			}),
		},
		{
			description: "dns custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[dnsCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.DNSCustomEndpoint = false
			}),
		},
		{
			description: "postgresql custom endpoint empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]bool) {
				flagValues[postgreSQLCustomEndpointFlag] = false
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.PostgreSQLCustomEndpoint = false
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}

			configureFlags(cmd)

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

			model := parseFlags(cmd)

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
