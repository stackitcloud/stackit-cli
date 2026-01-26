package create

import (
	"context"
	"strconv"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testRegion        = "eu01"
	testDisplayName   = "my-logs-instance"
	testDescription   = "my instance description"
	testAcl           = "198.51.100.14/24"
	testRetentionDays = 32
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &logs.APIClient{}
	testProjectId = uuid.NewString()
)

// Flags
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           testDisplayName,
		retentionDaysFlag:         strconv.Itoa(testRetentionDays),
		descriptionFlag:           testDescription,
		aclFlag:                   testAcl,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

// Input Model
func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		DisplayName:   utils.Ptr(testDisplayName),
		Description:   utils.Ptr(testDescription),
		RetentionDays: utils.Ptr(int64(testRetentionDays)),
		ACL:           utils.Ptr([]string{testAcl}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureRequest(mods ...func(request *logs.ApiCreateLogsInstanceRequest)) logs.ApiCreateLogsInstanceRequest {
	request := testClient.CreateLogsInstance(testCtx, testProjectId, testRegion)
	request = request.CreateLogsInstancePayload(logs.CreateLogsInstancePayload{
		DisplayName:   utils.Ptr(testDisplayName),
		Description:   utils.Ptr(testDescription),
		RetentionDays: utils.Ptr(int64(testRetentionDays)),
		Acl:           utils.Ptr([]string{testAcl}),
	})

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
			description: "optional flags omitted",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, descriptionFlag)
				delete(flagValues, aclFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.ACL = nil
			}),
		},
		{
			description: "no values provided",
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
			description: "display name missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "retention days missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, retentionDaysFlag)
			}),
			isValid: false,
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
		expectedRequest logs.ApiCreateLogsInstanceRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no optional values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.ACL = nil
			}),
			expectedRequest: fixtureRequest().CreateLogsInstancePayload(logs.CreateLogsInstancePayload{
				DisplayName:   utils.Ptr(testDisplayName),
				RetentionDays: utils.Ptr(int64(testRetentionDays)),
				Description:   nil,
				Acl:           nil,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(tt.expectedRequest, request,
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
		description string
		model       *inputModel
		instance    *logs.LogsInstance
		wantErr     bool
	}{
		{
			description: "nil response",
			instance:    nil,
			wantErr:     true,
		},
		{
			description: "model is nil",
			instance:    &logs.LogsInstance{},
			model:       nil,
			wantErr:     true,
		},
		{
			description: "global flag nil",
			instance:    &logs.LogsInstance{},
			model:       &inputModel{GlobalFlagModel: nil},
			wantErr:     true,
		},
		{
			description: "default output",
			instance:    &logs.LogsInstance{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			wantErr:     false,
		},
		{
			description: "json output",
			instance:    &logs.LogsInstance{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
			wantErr:     false,
		},
		{
			description: "yaml output",
			instance:    &logs.LogsInstance{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.YAMLOutputFormat}},
			wantErr:     false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(p, tt.model, "label", tt.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
