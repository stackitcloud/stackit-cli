package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

// Define a unique key for the context to avoid collisions
type testCtxKey struct{}

const (
	testRegion       = "eu01"
	testDisplayName  = "testuser"
	testPassword     = "Secret12345!"
	testUserType     = "intake"
	testDescription  = "This is a test user"
	testLabelsString = "env=test,team=dev"
)

var (
	// testCtx dummy context for testing purposes
	testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
	// testClient mock API client
	testClient    = &intake.APIClient{}
	testProjectId = uuid.NewString()
	testIntakeId  = uuid.NewString()

	testLabels = map[string]string{"env": "test", "team": "dev"}
)

// fixtureFlagValues generates a map of flag values for tests
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           testDisplayName,
		intakeIdFlag:              testIntakeId,
		passwordFlag:              testPassword,
		userTypeFlag:              testUserType,
		descriptionFlag:           testDescription,
		labelsFlag:                testLabelsString,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

// fixtureInputModel generates an input model for tests
func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		DisplayName: utils.Ptr(testDisplayName),
		IntakeId:    utils.Ptr(testIntakeId),
		Password:    utils.Ptr(testPassword),
		UserType:    utils.Ptr(testUserType),
		Description: utils.Ptr(testDescription),
		Labels:      utils.Ptr(testLabels),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// fixtureCreatePayload generates a CreateIntakeUserPayload for tests
func fixtureCreatePayload(mods ...func(payload *intake.CreateIntakeUserPayload)) intake.CreateIntakeUserPayload {
	userType := intake.UserType(testUserType)
	payload := intake.CreateIntakeUserPayload{
		DisplayName: utils.Ptr(testDisplayName),
		Password:    utils.Ptr(testPassword),
		Type:        &userType,
		Description: utils.Ptr(testDescription),
		Labels:      utils.Ptr(testLabels),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

// fixtureRequest generates an API request for tests
func fixtureRequest(mods ...func(request *intake.ApiCreateIntakeUserRequest)) intake.ApiCreateIntakeUserRequest {
	request := testClient.CreateIntakeUser(testCtx, testProjectId, testRegion, testIntakeId)
	request = request.CreateIntakeUserPayload(fixtureCreatePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			description: "intake id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, intakeIdFlag)
			}),
			isValid: false,
		},
		{
			description: "display name missing",
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
				displayNameFlag:           testDisplayName,
				intakeIdFlag:              testIntakeId,
				passwordFlag:              testPassword,
				userTypeFlag:              testUserType,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.Labels = nil
				// UserType has a default value in the command definition, so it should still be populated
				model.UserType = utils.Ptr(string(intake.USERTYPE_INTAKE))
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, func(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
				return parseInput(p, cmd)
			}, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest intake.ApiCreateIntakeUserRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no optionals",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.Labels = nil
				model.UserType = nil
			}),
			expectedRequest: fixtureRequest(func(request *intake.ApiCreateIntakeUserRequest) {
				*request = (*request).CreateIntakeUserPayload(fixtureCreatePayload(func(payload *intake.CreateIntakeUserPayload) {
					payload.Description = nil
					payload.Labels = nil
					payload.Type = nil
				}))
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(request, tt.expectedRequest,
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
	type args struct {
		model        *inputModel
		projectLabel string
		resp         *intake.IntakeUserResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default output",
			args: args{
				model:        fixtureInputModel(),
				projectLabel: "my-project",
				resp:         &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")},
			},
			wantErr: false,
		},
		{
			name: "default output - async",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Async = true
				}),
				projectLabel: "my-project",
				resp:         &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")},
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")},
			},
			wantErr: false,
		},
		{
			name: "nil response - default output",
			args: args{
				model: fixtureInputModel(),
				resp:  nil,
			},
			wantErr: false,
		},
		{
			name: "nil response - json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: nil,
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
