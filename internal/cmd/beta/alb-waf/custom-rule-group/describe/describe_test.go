package describe

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

type testCtxKey struct{}

var (
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "test")
	testProjectId    = uuid.NewString()
	testRegion       = "eu01"
	testClient       = &albwaf.APIClient{DefaultAPI: &albwaf.DefaultAPIService{}}
	testCustomRgName = "my-test-custom-rule-group"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testCustomRgName,
	}
	for _, m := range mods {
		m(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
	}
	for _, m := range mods {
		m(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
			Region:    testRegion,
		},
		Name: testCustomRgName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *albwaf.ApiGetCustomRuleGroupRequest)) albwaf.ApiGetCustomRuleGroupRequest {
	request := testClient.DefaultAPI.GetCustomRuleGroup(testCtx, testProjectId, testRegion, testCustomRgName)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argsValues    []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argsValues:    fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no arg values",
			argsValues:  []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "project id missing",
			argsValues:  fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argsValues:  fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argsValues:  fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argsValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description    string
		model          *inputModel
		expectedResult albwaf.ApiGetCustomRuleGroupRequest
	}{
		{
			description:    "base",
			model:          fixtureInputModel(),
			expectedResult: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedResult,
				cmp.AllowUnexported(tt.expectedResult, albwaf.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description  string
		outputFormat string
		crg          *albwaf.GetCustomRuleGroupResponse
		wantErr      bool
	}{
		{
			description:  "empty",
			outputFormat: "",
			crg:          nil,
			wantErr:      true,
		},
		{
			description:  "base",
			outputFormat: "",
			crg: &albwaf.GetCustomRuleGroupResponse{
				Name: testCustomRgName,
				Rules: []albwaf.GetCustomRule{
					{
						Id:          1,
						Description: utils.Ptr("block access to /admin"),
						Behavior: albwaf.GetBehavior{
							Action:   albwaf.ACTION_ACTION_DENY,
							Log:      true,
							LogMsg:   utils.Ptr("blocked"),
							Severity: albwaf.SEVERITY_SEVERITY_WARNING,
						},
						Conditions: []albwaf.Condition{
							{
								Operator: albwaf.ConditionOperator{
									Type:  albwaf.OPERATOR_OPERATOR_BEGINS_WITH,
									Value: utils.Ptr("/admin"),
								},
								Variable: albwaf.ConditionVariable{
									Type: albwaf.VARIABLE_VARIABLE_REQUEST_URI_RAW,
								},
								Transformations: []albwaf.Transformation{
									albwaf.TRANSFORMATION_TRANSFORMATION_LOWERCASE,
								},
							},
						},
					},
				},
				Usage: &albwaf.CRGUsage{
					Count: utils.Ptr(int32(1)),
					Items: []string{"my-waf"},
				},
			},
			wantErr: false,
		},
		{
			description:  "json output",
			outputFormat: print.JSONOutputFormat,
			crg: &albwaf.GetCustomRuleGroupResponse{
				Name:  testCustomRgName,
				Rules: []albwaf.GetCustomRule{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, tt.crg); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
