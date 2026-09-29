package create

import (
	"context"
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
	testTemplateId   = uuid.NewString()
	testAutomationId = uuid.NewString()
)

const (
	testRegion               = "eu01"
	testName                 = "my-automation"
	testDescription          = "my-description"
	testTriggerScheduleRrule = "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"
	testInputJson            = `{"kind": "VolumeRecoveryPointManagement"}`
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		templateIdFlag:            testTemplateId,
		nameFlag:                  testName,
		descriptionFlag:           testDescription,
		triggerScheduleRruleFlag:  testTriggerScheduleRrule,
		inputFlag:                 testInputJson,
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
		TemplateId:  testTemplateId,
		Name:        new(testName),
		Description: new(testDescription),
		Input: &automation.VolumeAutomationInput{
			Kind:                 "VolumeRecoveryPointManagement",
			AdditionalProperties: map[string]interface{}{},
		},
		TriggerScheduleRrule: new(testTriggerScheduleRrule),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *automation.ApiCreateVolumeAutomationRequest)) automation.ApiCreateVolumeAutomationRequest {
	payload := automation.CreateVolumeAutomationPayload{
		Name:        *automation.NewNullableString(new(testName)),
		Description: *automation.NewNullableString(new(testDescription)),
		Input: *automation.NewNullableVolumeAutomationInput(&automation.VolumeAutomationInput{
			Kind:                 "VolumeRecoveryPointManagement",
			AdditionalProperties: map[string]interface{}{},
		}),
		Triggers: *automation.NewNullableAutomationTriggers(&automation.AutomationTriggers{
			Schedule: *automation.NewNullableAutomationScheduleTrigger(&automation.AutomationScheduleTrigger{
				Rrule: "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR",
			}),
		}),
		TemplateId: testTemplateId,
	}
	request := testClient.DefaultAPI.CreateVolumeAutomation(testCtx, testProjectId, testRegion).CreateVolumeAutomationPayload(payload)
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
			description: "only required flags",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				templateIdFlag:            testTemplateId,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				TemplateId: testTemplateId,
			},
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "missing project id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "missing template id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, templateIdFlag)
			}),
			isValid: false,
		},
		{
			description: "template id is no uuid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[templateIdFlag] = "not-a-uuid"
			}),
			isValid: false,
		},
		{
			description: "invalid input json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[inputFlag] = "invalid-json"
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
		isValid         bool
		expectedRequest automation.ApiCreateVolumeAutomationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "only required fields",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				TemplateId: testTemplateId,
			},
			isValid: true,
			expectedRequest: fixtureRequest(func(request *automation.ApiCreateVolumeAutomationRequest) {
				payload := automation.CreateVolumeAutomationPayload{
					Name:        *automation.NewNullableString(nil),
					Description: *automation.NewNullableString(nil),
					Input:       *automation.NewNullableVolumeAutomationInput(nil),
					TemplateId:  testTemplateId,
					Triggers:    automation.NullableAutomationTriggers{},
				}
				*request = testClient.DefaultAPI.CreateVolumeAutomation(testCtx, testProjectId, testRegion).CreateVolumeAutomationPayload(payload)
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(
					tt.expectedRequest,
					automation.DefaultAPIService{},
					automation.NullableString{},
					automation.NullableVolumeAutomationInput{},
					automation.NullableAutomationTriggers{},
					automation.NullableAutomationScheduleTrigger{},
				),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match (-want, +got): %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		item         *automation.VolumeAutomation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil item",
			args: args{
				projectLabel: testProjectId,
				item:         nil,
			},
			wantErr: true,
		},
		{
			name: "valid response",
			args: args{
				projectLabel: testProjectId,
				item: &automation.VolumeAutomation{
					Id: testAutomationId,
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.projectLabel, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
