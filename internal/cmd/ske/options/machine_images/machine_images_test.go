package machine_images

import (
	"context"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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
		globalflags.RegionFlag: testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: globalflags.GlobalFlagModel{
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
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
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
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
		inputModel      *inputModel
		expectedRequest ske.ApiListProviderOptionsRequest
	}{
		{
			description:     "base",
			inputModel:      fixtureInputModel(),
			expectedRequest: testClient.ListProviderOptions(testCtx, testRegion),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, testClient, tt.inputModel)

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
		{
			name: "empty values",
			args: args{
				model: &inputModel{
					GlobalFlagModel: globalflags.GlobalFlagModel{},
				},
				options: &ske.ProviderOptions{
					MachineImages: &[]ske.MachineImage{},
				},
			},
			wantErr: false,
		},
		{
			name: "empty value in values",
			args: args{
				model: &inputModel{
					GlobalFlagModel: globalflags.GlobalFlagModel{},
				},
				options: &ske.ProviderOptions{
					MachineImages: &[]ske.MachineImage{{}},
				},
			},
			wantErr: false,
		},
		{
			name: "valid values",
			args: args{
				model: &inputModel{
					GlobalFlagModel: globalflags.GlobalFlagModel{},
				},
				options: &ske.ProviderOptions{
					MachineImages: &[]ske.MachineImage{
						{
							Name: utils.Ptr("image1"),
							Versions: &[]ske.MachineImageVersion{
								{
									Cri: &[]ske.CRI{
										{
											Name: ske.CRINAME_CONTAINERD.Ptr(),
										},
									},
									ExpirationDate: utils.Ptr(time.Now()),
									State:          utils.Ptr("supported"),
									Version:        utils.Ptr("0.00.0"),
								},
							},
						},
						{
							Name: utils.Ptr("zone2"),
						},
					},
				},
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
