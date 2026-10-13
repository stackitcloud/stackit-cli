package generatepayload

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var (
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient       = &automation.APIClient{DefaultAPI: &automation.DefaultAPIService{}}
	testProjectId    = uuid.NewString()
	testAutomationId = uuid.NewString()
)

const (
	testRegion   = "eu01"
	testFilePath = "./input-payload.json"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		automationIdFlag:          testAutomationId,
		filePathFlag:              testFilePath,
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
		AutomationId: testAutomationId,
		FilePath:     new(testFilePath),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *automation.ApiGetVolumeAutomationRequest)) automation.ApiGetVolumeAutomationRequest {
	request := testClient.DefaultAPI.GetVolumeAutomation(testCtx, testProjectId, testRegion, testAutomationId)
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
			description: "no file path",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, filePathFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.FilePath = nil
			}),
		},
		{
			description:   "no values",
			flagValues:    map[string]string{},
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "missing automation id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, automationIdFlag)
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "automation id is no uuid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[automationIdFlag] = "not-a-uuid"
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "missing project id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid:       false,
			expectedModel: nil,
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
		isValid         bool
		expectedRequest automation.ApiGetVolumeAutomationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest, automation.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match (-want, +got): %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tempDir := t.TempDir()
	type args struct {
		filePath *string
		item     *automation.VolumeAutomation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil item",
			args: args{
				filePath: nil,
				item:     nil,
			},
			wantErr: true,
		},
		{
			name: "terminal output",
			args: args{
				filePath: nil,
				item: &automation.VolumeAutomation{
					Id: testAutomationId,
					Input: *automation.NewNullableVolumeAutomationInput(&automation.VolumeAutomationInput{
						Kind: "VolumeRecoveryPointManagement",
					}),
				},
			},
			wantErr: false,
		},
		{
			name: "write to file",
			args: args{
				filePath: new(filepath.Join(tempDir, "payload.json")),
				item: &automation.VolumeAutomation{
					Id: testAutomationId,
					Input: *automation.NewNullableVolumeAutomationInput(&automation.VolumeAutomationInput{
						Kind: "VolumeRecoveryPointManagement",
					}),
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.filePath, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.filePath != nil && !tt.wantErr {
				if _, err := os.Stat(*tt.args.filePath); os.IsNotExist(err) {
					t.Errorf("expected file %s to be created", *tt.args.filePath)
				}
			}
		})
	}
}
