package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &mongodbflex.APIClient{}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:  testProjectId,
		instanceIdFlag: testInstanceId,
		usernameFlag:   "johndoe",
		databaseFlag:   "default",
		roleFlag:       "read",
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
		},
		InstanceId: testInstanceId,
		Username:   utils.Ptr("johndoe"),
		Database:   utils.Ptr("default"),
		Roles:      utils.Ptr([]string{"read"}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *mongodbflex.ApiCreateUserRequest)) mongodbflex.ApiCreateUserRequest {
	request := testClient.CreateUser(testCtx, testProjectId, testInstanceId)
	request = request.CreateUserPayload(mongodbflex.CreateUserPayload{
		Username: utils.Ptr("johndoe"),
		Database: utils.Ptr("default"),
		Roles:    utils.Ptr([]string{"read"}),
	})

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
			description: "no username specified",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, usernameFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Username = nil
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "database missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, databaseFlag)
			}),
			isValid: false,
		},
		{
			description: "roles missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, roleFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Roles = &rolesDefault
			}),
		},
		{
			description: "invalid role",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[roleFlag] = "invalid-role"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			configureFlags(cmd)

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

			model, err := parseInput(cmd)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest mongodbflex.ApiCreateUserRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no username specified",
			model: fixtureInputModel(func(model *inputModel) {
				model.Username = nil
			}),
			expectedRequest: fixtureRequest().CreateUserPayload(mongodbflex.CreateUserPayload{
				Database: utils.Ptr("default"),
				Roles:    utils.Ptr([]string{"read"}),
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
