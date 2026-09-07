package list

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

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")

var testClient = &modelexperiments.APIClient{
	DefaultAPI: modelexperiments.DefaultAPIServiceMock{},
}

var testProjectId = uuid.NewString()

func fixtureFlagValues(
	mods ...func(*map[string]string),
) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: "eu01",
	}

	for _, mod := range mods {
		mod(&flagValues)
	}

	return flagValues
}

func fixtureInputModel(
	mods ...func(*inputModel),
) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    "eu01",
			Verbosity: globalflags.VerbosityDefault,
		},
		Region: "eu01",
	}

	for _, mod := range mods {
		mod(model)
	}

	return model
}

func fixtureRequest(
	mods ...func(*modelexperiments.ApiListInstancesRequest),
) modelexperiments.ApiListInstancesRequest {
	request := testClient.DefaultAPI.ListInstances(
		testCtx,
		testProjectId,
		"eu01",
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
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues *map[string]string) {
				delete(*flagValues, projectIdFlag)
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

func TestBuildListInstancesRequest(t *testing.T) {
	request := buildListInstancesRequest(
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
		name    string
		resp    *modelexperiments.ListInstancesResponse
		wantErr bool
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name: "instances",
			resp: &modelexperiments.ListInstancesResponse{
				Instances: []modelexperiments.Instance{},
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
