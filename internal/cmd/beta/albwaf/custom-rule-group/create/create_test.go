package create

import (
	"context"
	"testing"

	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

const (
	testRegion = "eu02"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &albwaf.APIClient{DefaultAPI: &albwaf.DefaultAPIService{}}
	testProjectId = uuid.NewString()
)

var testPayload = &albwaf.CreateCustomRuleGroupPayload{
	Name:  "test-custom-rule-group",
	Rules: []albwaf.CreateCustomRule{},
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		payloadFlag: `{
  "name": "test-custom-rule-group",
  "rules": []
}`,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		Payload: testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *albwaf.ApiCreateCustomRuleGroupRequest)) albwaf.ApiCreateCustomRuleGroupRequest {
	request := testClient.DefaultAPI.CreateCustomRuleGroup(testCtx, testProjectId, testRegion)
	request = request.CreateCustomRuleGroupPayload(*testPayload)
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
			description: "payload is missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, payloadFlag)
			}),
			isValid: false,
		},
		{
			description: "payload is empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid:       false,
			expectedModel: fixtureInputModel(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithOptions(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, nil, tt.isValid, []testutils.TestingOption{
				testutils.WithCmpOptions(cmpopts.EquateEmpty()),
			})
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest albwaf.ApiCreateCustomRuleGroupRequest
		isValid         bool
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

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, albwaf.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description  string
		outputFormat string
		resp         *albwaf.GetCustomRuleGroupResponse
		wantErr      bool
	}{
		{
			description:  "empty",
			outputFormat: "",
			resp:         nil,
			wantErr:      false,
		},
		{
			description:  "base",
			outputFormat: "",
			resp: &albwaf.GetCustomRuleGroupResponse{
				Name: "test-custom-rule-group",
			},
			wantErr: false,
		},
		{
			description:  "json output",
			outputFormat: print.JSONOutputFormat,
			resp: &albwaf.GetCustomRuleGroupResponse{
				Name: "test-custom-rule-group",
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
