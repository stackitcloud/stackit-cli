package list

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

const testRegion = "eu01"

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
	testProjectId = uuid.NewString()
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
	}

	for _, mod := range mods {
		mod(model)
	}

	return model
}

func fixtureRequest(mods ...func(*modelexperiments.ApiListInstancesRequest)) modelexperiments.ApiListInstancesRequest {
	request := testClient.DefaultAPI.ListInstances(
		testCtx,
		testProjectId,
		testRegion,
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
			description: "limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "10"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(10))
			}),
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "limit must be greater than zero",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
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

func TestApplyLimit(t *testing.T) {
	tests := []struct {
		name      string
		instances []modelexperiments.Instance
		limit     *int64
		nilResp   bool
		wantCount int
	}{
		{
			name:      "no limit",
			instances: make([]modelexperiments.Instance, 3),
			wantCount: 3,
		},
		{
			name:      "limit below result count",
			instances: make([]modelexperiments.Instance, 3),
			limit:     utils.Ptr(int64(2)),
			wantCount: 2,
		},
		{
			name:      "limit above result count",
			instances: make([]modelexperiments.Instance, 2),
			limit:     utils.Ptr(int64(3)),
			wantCount: 2,
		},
		{
			name:      "nil response",
			limit:     utils.Ptr(int64(1)),
			nilResp:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &modelexperiments.ListInstancesResponse{Instances: tt.instances}
			if tt.nilResp {
				resp = nil
			}

			applyLimit(resp, tt.limit)

			if resp == nil {
				if !tt.nilResp {
					t.Fatal("applyLimit() unexpectedly returned a nil response")
				}
				return
			}
			if tt.nilResp {
				t.Fatal("applyLimit() changed a nil response")
			}
			if got := len(resp.Instances); got != tt.wantCount {
				t.Errorf("applyLimit() count = %d, want %d", got, tt.wantCount)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name       string
		resp       *modelexperiments.ListInstancesResponse
		wantErr    bool
		wantOutput []string
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name: "instances",
			resp: &modelexperiments.ListInstancesResponse{
				Instances: []modelexperiments.Instance{
					{
						Id:         "instance-id",
						Name:       "instance-name",
						BucketName: utils.Ptr("model-experiments-bucket"),
						Region:     utils.Ptr(testRegion),
						State:      "READY",
					},
				},
			},
			wantErr:    false,
			wantOutput: []string{"ID", "NAME", "instance-id", "instance-name", "model-experiments-bucket", "READY"},
		},
		{
			name: "empty instances",
			resp: &modelexperiments.ListInstancesResponse{
				Instances: []modelexperiments.Instance{},
			},
			wantErr:    false,
			wantOutput: []string{`No instances found for project "test-project"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()
			if err := outputResult(
				params.Printer,
				"",
				"test-project",
				tt.resp,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"outputResult() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}

			for _, expected := range tt.wantOutput {
				if !strings.Contains(params.Out.String(), expected) {
					t.Errorf("outputResult() output = %q, want it to contain %q", params.Out.String(), expected)
				}
			}
		})
	}
}
