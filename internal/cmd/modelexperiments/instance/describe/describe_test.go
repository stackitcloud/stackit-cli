package describe

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

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testInstanceId,
	}

	for _, mod := range mods {
		mod(argValues)
	}

	return argValues
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

func fixtureRequest(mods ...func(*modelexperiments.ApiGetInstanceRequest)) modelexperiments.ApiGetInstanceRequest {
	request := testClient.DefaultAPI.GetInstance(
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
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
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
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

func TestBuildDescribeInstanceRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest modelexperiments.ApiGetInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(
				request,
				tt.expectedRequest,
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
	tests := []struct {
		name    string
		resp    *modelexperiments.GetInstanceResponse
		wantErr bool
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name: "empty instance",
			resp: &modelexperiments.GetInstanceResponse{
				Instance: modelexperiments.Instance{},
			},
			wantErr: false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(
				params.Printer,
				"",
				tt.resp,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"outputResult() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
