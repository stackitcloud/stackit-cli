package set

import (
	"os"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
)

const testProfile = "test-profile"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testProfile,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		Profile: testProfile,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		isValid       bool
		profileEnvVar *string
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argValues:     fixtureArgValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			isValid:     false,
		},
		{
			description: "some global flag",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.VerbosityFlag: globalflags.DebugVerbosity,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.GlobalFlagModel.Verbosity = globalflags.DebugVerbosity
			}),
		},
		{
			description: "invalid profile",
			argValues:   []string{"invalid-profile-&"},
			isValid:     false,
		},
		{
			description:   "profile from env",
			argValues:     []string{},
			profileEnvVar: utils.Ptr(testProfile),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description:   "profile from env, but empty",
			argValues:     []string{},
			profileEnvVar: utils.Ptr(""),
			isValid:       false,
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

			// Set the profile env var or unset existing
			oldProfileEnv := os.Getenv(config.ProfileEnvVar)
			if tt.profileEnvVar != nil {
				err = os.Setenv(config.ProfileEnvVar, *tt.profileEnvVar)
				if err != nil {
					t.Fatalf("setting profile env var: %v", err)
				}
			} else {
				err = os.Unsetenv(config.ProfileEnvVar)
				if err != nil {
					t.Fatalf("unsetting profile env var: %v", err)
				}
			}

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}

			// Reset the env state
			err = os.Setenv(config.ProfileEnvVar, oldProfileEnv)
			if err != nil {
				t.Fatalf("setting back profile env var: %v", err)
			}
		})
	}
}
