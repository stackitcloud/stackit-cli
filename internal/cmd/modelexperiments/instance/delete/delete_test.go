package delete

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

const testRegion = "eu01"

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
	}

	for _, mod := range mods {
		mod(flagValues)
	}

	return flagValues
}

func fixtureInputModel(mods ...func(*inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId: testInstanceId,
	}

	for _, mod := range mods {
		mod(model)
	}

	return model
}

func fixtureRequest(mods ...func(*modelexperiments.ApiDeleteInstanceRequest)) modelexperiments.ApiDeleteInstanceRequest {
	request := testClient.DefaultAPI.DeleteInstance(
		testCtx,
		testProjectId,
		testRegion,
		testInstanceId,
	)

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
			argValues:     []string{testInstanceId},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "project id missing",
			argValues:   []string{testInstanceId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(
				t,
				NewCmd,
				parseInput,
				tt.expectedModel,
				tt.argValues,
				tt.flagValues,
				tt.isValid,
			)
		})
	}
}

func TestBuildDeleteInstanceRequest(t *testing.T) {
	request := buildRequest(
		testCtx,
		fixtureInputModel(),
		testClient,
	)

	expectedRequest := fixtureRequest()

	diff := cmp.Diff(
		request,
		expectedRequest,
		cmp.AllowUnexported(expectedRequest),
		cmpopts.EquateComparable(testCtx),
	)

	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name       string
		resp       *modelexperiments.DeleteInstanceResponse
		async      bool
		wantErr    bool
		wantOutput string
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name:       "response",
			resp:       &modelexperiments.DeleteInstanceResponse{},
			wantOutput: "Deleted AI Model Experiments instance.\n",
		},
		{
			name:       "async response",
			resp:       &modelexperiments.DeleteInstanceResponse{},
			async:      true,
			wantOutput: "Triggered deletion of AI Model Experiments instance.\n",
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params.Out.Reset()
			if err := outputResult(
				params.Printer,
				"",
				tt.async,
				tt.resp,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"outputResult() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			if got := params.Out.String(); got != tt.wantOutput {
				t.Errorf("outputResult() output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}
