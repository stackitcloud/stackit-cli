package describe

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "test")
	testProjectId = uuid.NewString()
	testRegion    = "eu01"
	testClient    = &albwaf.APIClient{DefaultAPI: &albwaf.DefaultAPIService{}}
	testName      = "my-managed-rule-set"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testName,
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
		Name: testName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *albwaf.ApiGetManagedRuleSetRequest)) albwaf.ApiGetManagedRuleSetRequest {
	request := testClient.DefaultAPI.GetManagedRuleSet(testCtx, testProjectId, testRegion, testName)
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
		expectedRequest albwaf.ApiGetManagedRuleSetRequest
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
				cmp.AllowUnexported(tt.expectedRequest, albwaf.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		resp         *albwaf.GetManagedRuleSetResponse
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
			name: "base",
			args: args{
				resp: &albwaf.GetManagedRuleSetResponse{
					Name:    testName,
					Type:    albwaf.TYPE_TYPE_OWASP_CRS,
					Version: "v1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "with groups",
			args: args{
				resp: &albwaf.GetManagedRuleSetResponse{
					Name:    testName,
					Type:    albwaf.TYPE_TYPE_OWASP_CRS,
					Version: "v1.0.0",
					Groups: &map[string]albwaf.MRSRuleGroup{
						"942": {
							Description: "SQL Injection",
							GroupName:   "SQL Injection",
							Rules: &map[string]albwaf.MRSRule{
								"942100": {
									Description: "SQL Injection Attack Detected",
									Mode:        albwaf.MODE_MODE_ENABLED,
									Severity:    "CRITICAL",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
