package update

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

type testCtxKey struct{}

const (
	testRegion = "eu01"
)

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &intake.APIClient{}
	testProjectId = uuid.NewString()
	testIntakeId  = uuid.NewString()
	testUserId    = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{testUserId}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		intakeIdFlag:              testIntakeId,
		displayNameFlag:           "new-display-name",
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
		IntakeId:    testIntakeId,
		UserId:      testUserId,
		DisplayName: utils.Ptr("new-display-name"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *intake.ApiUpdateIntakeUserRequest)) intake.ApiUpdateIntakeUserRequest {
	request := testClient.UpdateIntakeUser(testCtx, testProjectId, testRegion, testIntakeId, testUserId)
	payload := intake.UpdateIntakeUserPayload{
		DisplayName: utils.Ptr("new-display-name"),
	}
	request = request.UpdateIntakeUserPayload(payload)
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
			description: "no optional flags provided",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				intakeIdFlag:              testIntakeId,
			},
			isValid: false,
		},
		{
			description: "update all fields",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[descriptionFlag] = "new description"
				flagValues[labelsFlag] = "env=prod,team=sre"
				flagValues[userTypeFlag] = "dead-letter"
				flagValues[passwordFlag] = "NewSecret123!"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("new description")
				model.Labels = utils.Ptr(map[string]string{"env": "prod", "team": "sre"})
				model.UserType = utils.Ptr("dead-letter")
				model.Password = utils.Ptr("NewSecret123!")
			}),
		},
		{
			description: "no args",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "intake-id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, intakeIdFlag)
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
		description string
		model       *inputModel
		expectedReq intake.ApiUpdateIntakeUserRequest
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			expectedReq: fixtureRequest(),
		},
		{
			description: "update description",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Description = utils.Ptr("new-desc")
			}),
			expectedReq: fixtureRequest(func(request *intake.ApiUpdateIntakeUserRequest) {
				payload := intake.UpdateIntakeUserPayload{
					Description: utils.Ptr("new-desc"),
				}
				*request = (*request).UpdateIntakeUserPayload(payload)
			}),
		},
		{
			description: "update all fields",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = utils.Ptr("another-name")
				model.Description = utils.Ptr("final-desc")
				model.Labels = utils.Ptr(map[string]string{"a": "b"})
				model.UserType = utils.Ptr("dead-letter")
				model.Password = utils.Ptr("Secret123!")
			}),
			expectedReq: fixtureRequest(func(request *intake.ApiUpdateIntakeUserRequest) {
				userType := intake.UserType("dead-letter")
				payload := intake.UpdateIntakeUserPayload{
					DisplayName: utils.Ptr("another-name"),
					Description: utils.Ptr("final-desc"),
					Labels:      utils.Ptr(map[string]string{"a": "b"}),
					Type:        &userType,
					Password:    utils.Ptr("Secret123!"),
				}
				*request = (*request).UpdateIntakeUserPayload(payload)
			}),
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
		outputFormat string
		projectLabel string
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
			args:    args{outputFormat: "default", projectLabel: "my-project", intakeId: "intake-id-123", resp: &intake.IntakeUserResponse{}},
			wantErr: false,
		},
		{
			name:    "json output",
			args:    args{outputFormat: print.JSONOutputFormat, resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")}},
			wantErr: false,
		},
		{
			name:    "nil response",
			args:    args{outputFormat: print.JSONOutputFormat, resp: nil},
			wantErr: false,
		},
		{
			name:    "nil response - default output",
			args:    args{outputFormat: "default", resp: nil},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: tt.args.outputFormat}, IntakeId: tt.args.intakeId}, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
