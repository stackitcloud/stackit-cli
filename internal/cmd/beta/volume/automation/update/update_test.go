package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient       = &automation.APIClient{DefaultAPI: &automation.DefaultAPIService{}}
	testProjectId    = uuid.NewString()
	testAutomationId = uuid.NewString()
)

const (
	testRegion               = "eu01"
	testName                 = "my-automation"
	testDescription          = "my-description"
	testTriggerScheduleRrule = "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"
	testInputJson            = `{"kind": "VolumeRecoveryPointManagement"}`
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testAutomationId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
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
		AutomationId: testAutomationId,
		Name:         new(testName),
		Description:  new(testDescription),
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

func fixtureRequest(mods ...func(request *automation.ApiPartialUpdateVolumeAutomationRequest)) automation.ApiPartialUpdateVolumeAutomationRequest {
	payload := automation.PartialUpdateVolumeAutomationPayload{
		Name:        *automation.NewNullableString(new(testName)),
		Description: *automation.NewNullableString(new(testDescription)),
		Input: *automation.NewNullableVolumeAutomationInput(&automation.VolumeAutomationInput{
			Kind:                 "VolumeRecoveryPointManagement",
			AdditionalProperties: map[string]interface{}{},
		}),
		Triggers: *automation.NewNullableAutomationTriggers(&automation.AutomationTriggers{
			Schedule: *automation.NewNullableAutomationScheduleTrigger(&automation.AutomationScheduleTrigger{
				Rrule: testTriggerScheduleRrule,
			}),
		}),
	}
	request := testClient.DefaultAPI.PartialUpdateVolumeAutomation(testCtx, testProjectId, testRegion, testAutomationId).PartialUpdateVolumeAutomationPayload(payload)
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
			description: "disable trigger schedule",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, triggerScheduleRruleFlag)
				flagValues[disableTriggerSchedule] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.TriggerScheduleRrule = nil
				model.DisableTriggerSchedule = utils.Ptr(true)
			}),
		},
		{
			description: "only name flag",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				nameFlag:                  testName,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				AutomationId: testAutomationId,
				Name:         new(testName),
			},
		},
		{
			description: "no values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no update flags",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
			},
			isValid: false,
		},
		{
			description: "mutually exclusive trigger flags",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[disableTriggerSchedule] = "true"
			}),
			isValid: false,
		},
		{
			description: "missing project id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "missing automation id arg",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "automation id is no uuid",
			argValues:   []string{"not-a-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid input json",
			argValues:   fixtureArgValues(),
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
		expectedRequest automation.ApiPartialUpdateVolumeAutomationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "disable trigger schedule",
			model: fixtureInputModel(func(model *inputModel) {
				model.TriggerScheduleRrule = nil
				model.DisableTriggerSchedule = utils.Ptr(true)
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *automation.ApiPartialUpdateVolumeAutomationRequest) {
				payload := automation.PartialUpdateVolumeAutomationPayload{
					Name:        *automation.NewNullableString(new(testName)),
					Description: *automation.NewNullableString(new(testDescription)),
					Input: *automation.NewNullableVolumeAutomationInput(&automation.VolumeAutomationInput{
						Kind:                 "VolumeRecoveryPointManagement",
						AdditionalProperties: map[string]interface{}{},
					}),
					Triggers: *automation.NewNullableAutomationTriggers(&automation.AutomationTriggers{
						Schedule: *automation.NewNullableAutomationScheduleTrigger(nil),
					}),
				}
				*request = testClient.DefaultAPI.PartialUpdateVolumeAutomation(testCtx, testProjectId, testRegion, testAutomationId).PartialUpdateVolumeAutomationPayload(payload)
			}),
		},
		{
			description: "only name field",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				AutomationId: testAutomationId,
				Name:         new(testName),
			},
			isValid: true,
			expectedRequest: fixtureRequest(func(request *automation.ApiPartialUpdateVolumeAutomationRequest) {
				payload := automation.PartialUpdateVolumeAutomationPayload{
					Name: *automation.NewNullableString(new(testName)),
				}
				*request = testClient.DefaultAPI.PartialUpdateVolumeAutomation(testCtx, testProjectId, testRegion, testAutomationId).PartialUpdateVolumeAutomationPayload(payload)
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
