package list

import (
	"context"
	"strconv"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var (
	testCtx             = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient          = &iaas.APIClient{}
	testProjectId       = uuid.NewString()
	testLabels          = "fooKey=fooValue,barKey=barValue,bazKey=bazValue"
	testLimit     int64 = 10
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:     testProjectId,
		labelSelectorFlag: testLabels,
		limitFlag:         strconv.Itoa(int(testLimit)),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectId, Verbosity: globalflags.VerbosityDefault},
		LabelSelector:   utils.Ptr(testLabels),
		Limit:           &testLimit,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiListImagesRequest)) iaas.ApiListImagesRequest {
	request := testClient.ListImages(testCtx, testProjectId)
	request = request.LabelSelector(testLabels)
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "no labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelSelectorFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.LabelSelector = nil
			}),
		},
		{
			description: "single label",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelSelectorFlag] = "foo=bar"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.LabelSelector = utils.Ptr("foo=bar")
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
		model           *inputModel
		expectedRequest iaas.ApiListImagesRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},

		{
			description: "no labels",
			model: fixtureInputModel(func(model *inputModel) {
				model.LabelSelector = utils.Ptr("")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiListImagesRequest) {
				*request = (*request).LabelSelector("")
			}),
		},
		{
			description: "single label",
			model: fixtureInputModel(func(model *inputModel) {
				model.LabelSelector = utils.Ptr("foo=bar")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiListImagesRequest) {
				*request = (*request).LabelSelector("foo=bar")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
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

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat string
		items        []iaas.Image
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				outputFormat: "",
				items:        []iaas.Image{},
			},
			wantErr: false,
		},
		{
			name: "output format json",
			args: args{
				outputFormat: print.JSONOutputFormat,
				items:        []iaas.Image{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.items); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
