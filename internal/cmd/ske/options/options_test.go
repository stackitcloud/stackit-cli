package options

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{}

const testRegion = "eu01"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		availabilityZonesFlag:  "false",
		kubernetesVersionsFlag: "false",
		machineImagesFlag:      "false",
		machineTypesFlag:       "false",
		volumeTypesFlag:        "false",
		globalflags.RegionFlag: testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModelAllFalse(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel:    globalflags.GlobalFlagModel{Region: testRegion, Verbosity: globalflags.VerbosityDefault},
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
		GlobalFlagModel:    globalflags.GlobalFlagModel{Region: testRegion, Verbosity: globalflags.VerbosityDefault},
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
		argValues     []string
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: fixtureInputModelAllTrue(func(model *inputModel) {
				model.Region = ""
			}),
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
				model.Region = ""
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
				model.Region = ""
			}),
		},
		{
			description: "some values 3",
			flagValues: map[string]string{
				kubernetesVersionsFlag: "false",
				machineTypesFlag:       "false",
			},
			isValid: true,
			expectedModel: fixtureInputModelAllTrue(func(model *inputModel) {
				model.Region = ""
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
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
			expectedRequest: testClient.ListProviderOptions(testCtx, testRegion),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, testClient, fixtureInputModelAllTrue())

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

func TestOutputResult(t *testing.T) {
	type args struct {
		model   *inputModel
		options *ske.ProviderOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "missing model",
			args: args{
				options: &ske.ProviderOptions{},
			},
			wantErr: true,
		},
		{
			name: "missing options",
			args: args{
				model: &inputModel{
					GlobalFlagModel: globalflags.GlobalFlagModel{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty input model",
			args: args{
				model:   &inputModel{},
				options: &ske.ProviderOptions{},
			},
			wantErr: false,
		},
		{
			name: "set model and options",
			args: args{
				model: &inputModel{
					GlobalFlagModel: globalflags.GlobalFlagModel{},
				},
				options: &ske.ProviderOptions{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutputResultAsTable(t *testing.T) {
	type args struct {
		options *ske.ProviderOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "empty options",
			args: args{
				options: &ske.ProviderOptions{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResultAsTable(p, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("outputResultAsTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
