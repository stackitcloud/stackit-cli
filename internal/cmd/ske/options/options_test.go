package options

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		availabilityZonesFlag:  "false",
		kubernetesVersionsFlag: "false",
		machineImagesFlag:      "false",
		machineTypesFlag:       "false",
		volumeTypesFlag:        "false",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModelAllFalse(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel:    &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		AvailabilityZones:  false,
		KubernetesVersions: false,
		MachineImages:      false,
		MachineTypes:       false,
		VolumeTypes:        false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureInputModelAllTrue(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel:    &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		AvailabilityZones:  true,
		KubernetesVersions: true,
		MachineImages:      true,
		MachineTypes:       true,
		VolumeTypes:        true,
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
			expectedModel: fixtureInputModelAllTrue(),
		},
		{
			description:   "no values",
			flagValues:    map[string]string{},
			isValid:       true,
			expectedModel: fixtureInputModelAllTrue(),
		},
		{
			description: "some values 1",
			flagValues: map[string]string{
				availabilityZonesFlag:  "true",
				kubernetesVersionsFlag: "false",
			},
			isValid: true,
			expectedModel: fixtureInputModelAllFalse(func(model *inputModel) {
				model.AvailabilityZones = true
			}),
		},
		{
			description: "some values 2",
			flagValues: map[string]string{
				kubernetesVersionsFlag: "true",
				machineImagesFlag:      "false",
				machineTypesFlag:       "true",
			},
			isValid: true,
			expectedModel: fixtureInputModelAllFalse(func(model *inputModel) {
				model.KubernetesVersions = true
				model.MachineTypes = true
			}),
		},
		{
			description: "some values 3",
			flagValues: map[string]string{
				kubernetesVersionsFlag: "false",
				machineTypesFlag:       "false",
			},
			isValid:       true,
			expectedModel: fixtureInputModelAllTrue(),
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

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd)
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

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		expectedRequest ske.ApiListProviderOptionsRequest
	}{
		{
			description:     "base",
			expectedRequest: testClient.ListProviderOptions(testCtx),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
