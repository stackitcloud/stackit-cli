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

type testContextKey struct{}

const testRegion = "eu01"

var (
	testContext    = context.WithValue(context.Background(), testContextKey{}, "token-list")
	testClient     = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
	testProjectID  = uuid.NewString()
	testInstanceID = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,
		instanceIDFlag:            testInstanceID,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectID,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceID: testInstanceID,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest() modelexperiments.ApiListInstanceTokensRequest {
	return testClient.DefaultAPI.ListInstanceTokens(testContext, testProjectID, testRegion, testInstanceID)
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
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIDFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "instance id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "limit 0",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
		{
			description: "limit negative",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "-1"
			}),
			isValid: false,
		},
		{
			description: "limit invalid string",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "invalid"
			}),
			isValid: false,
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	request := buildRequest(testContext, fixtureInputModel(), testClient)
	expectedRequest := fixtureRequest()
	diff := cmp.Diff(request, expectedRequest, cmp.AllowUnexported(expectedRequest), cmpopts.EquateComparable(testContext))
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestApplyLimit(t *testing.T) {
	tests := []struct {
		name      string
		tokens    []modelexperiments.TokenMetadata
		limit     *int64
		nilResp   bool
		wantCount int
	}{
		{
			name:      "no limit",
			tokens:    make([]modelexperiments.TokenMetadata, 3),
			wantCount: 3,
		},
		{
			name:      "limit below result count",
			tokens:    make([]modelexperiments.TokenMetadata, 3),
			limit:     utils.Ptr(int64(2)),
			wantCount: 2,
		},
		{
			name:      "limit above result count",
			tokens:    make([]modelexperiments.TokenMetadata, 2),
			limit:     utils.Ptr(int64(3)),
			wantCount: 2,
		},
		{
			name:    "nil response",
			limit:   utils.Ptr(int64(1)),
			nilResp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &modelexperiments.ListInstanceTokensResponse{Tokens: tt.tokens}
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
			if got := len(resp.Tokens); got != tt.wantCount {
				t.Errorf("applyLimit() count = %d, want %d", got, tt.wantCount)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name       string
		resp       *modelexperiments.ListInstanceTokensResponse
		wantErr    bool
		wantOutput []string
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name: "tokens",
			resp: &modelexperiments.ListInstanceTokensResponse{
				Tokens: []modelexperiments.TokenMetadata{
					{
						Id:     "token-id",
						Name:   "token-name",
						Region: testRegion,
						State:  "READY",
					},
				},
			},
			wantOutput: []string{"ID", "NAME", "token-id", "token-name", testRegion, "READY"},
		},
		{
			name: "empty tokens",
			resp: &modelexperiments.ListInstanceTokensResponse{
				Tokens: []modelexperiments.TokenMetadata{},
			},
			wantOutput: []string{`No tokens found for instance "test-instance" in project "test-project"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()
			if err := outputResult(params.Printer, fixtureInputModel(), "test-project", "test-instance", tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, expected := range tt.wantOutput {
				if !strings.Contains(params.Out.String(), expected) {
					t.Errorf("outputResult() output = %q, want it to contain %q", params.Out.String(), expected)
				}
			}
		})
	}
}
