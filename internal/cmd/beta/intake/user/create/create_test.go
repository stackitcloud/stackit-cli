package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

type testCtxKey struct{}

const (
	testRegion = "eu01"

	testDisplayName  = "testuser"
	testPassword     = "my-secret-password"
	testDescription  = "This is a test user"
	testType         = "intake"
	testLabelsString = "env=test,owner=team-blue"
)

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &intake.APIClient{}
	testProjectId = uuid.NewString()
	testIntakeId  = uuid.NewString()

	testLabels = map[string]string{"env": "test", "owner": "team-blue"}
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		intakeIdFlag:              testIntakeId,
		displayNameFlag:           testDisplayName,
		passwordFlag:              testPassword,
		descriptionFlag:           testDescription,
		typeFlag:                  testType,
		labelsFlag:                testLabelsString,
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
		IntakeId:    utils.Ptr(testIntakeId),
		DisplayName: utils.Ptr(testDisplayName),
		Password:    utils.Ptr(testPassword),
		Description: utils.Ptr(testDescription),
		Type:        utils.Ptr(testType),
		Labels:      utils.Ptr(testLabels),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureCreatePayload(mods ...func(payload *intake.CreateIntakeUserPayload)) intake.CreateIntakeUserPayload {
	payload := intake.CreateIntakeUserPayload{
		DisplayName: utils.Ptr(testDisplayName),
		Password:    utils.Ptr(testPassword),
		Description: utils.Ptr(testDescription),
		Type:        (*intake.UserType)(utils.Ptr(testType)),
		Labels:      utils.Ptr(testLabels),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
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
			description: "intake-id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, intakeIdFlag)
			}),
			isValid: false,
		},
		{
			description: "display-name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "password missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, passwordFlag)
			}),
			isValid: false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				intakeIdFlag:              testIntakeId,
				displayNameFlag:           testDisplayName,
				passwordFlag:              testPassword,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.Type = nil
				model.Labels = nil
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(tt.expectedModel, model)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expectedReq intake.ApiCreateIntakeUserRequest
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			expectedReq: testClient.CreateIntakeUser(testCtx, testProjectId, testRegion, testIntakeId).
				CreateIntakeUserPayload(fixtureCreatePayload()),
		},
		{
			description: "no optionals",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.Type = nil
				model.Labels = nil
			}),
			expectedReq: testClient.CreateIntakeUser(testCtx, testProjectId, testRegion, testIntakeId).
				CreateIntakeUserPayload(fixtureCreatePayload(func(payload *intake.CreateIntakeUserPayload) {
					payload.Description = nil
					payload.Type = nil
					payload.Labels = nil
				})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(tt.expectedReq, request,
				cmp.AllowUnexported(request),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		model        *inputModel
		outputFormat string
		intakeId     string
		resp         *intake.IntakeUserResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "default output",
			args:    args{outputFormat: "default", intakeId: "intake-id-123", resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")}, model: fixtureInputModel()},
			wantErr: false,
		},
		{
			name:    "json output",
			args:    args{outputFormat: print.JSONOutputFormat, resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")}, model: fixtureInputModel()},
			wantErr: false,
		},
		{
			name: "yaml output",
			args: args{outputFormat: print.YAMLOutputFormat, resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")}, model: fixtureInputModel(func(model *inputModel) {
				model.OutputFormat = print.JSONOutputFormat
			})},
			wantErr: false,
		},
		{
			name: "nil response - json output",
			args: args{outputFormat: print.JSONOutputFormat, resp: nil, model: fixtureInputModel(func(model *inputModel) {
				model.OutputFormat = print.JSONOutputFormat
			})},
			wantErr: false,
		},
		{
			name:    "nil response - default output",
			args:    args{outputFormat: "default", intakeId: "intake-id-123", resp: nil, model: fixtureInputModel()},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.intakeId, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
