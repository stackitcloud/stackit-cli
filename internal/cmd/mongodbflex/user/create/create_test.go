package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	testRegion = "eu02"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &mongodbflex.APIClient{}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceIdFlag:            testInstanceId,
		usernameFlag:              "johndoe",
		databaseFlag:              "default",
		roleFlag:                  "read",
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
	request := testClient.CreateUser(testCtx, testProjectId, testInstanceId, testRegion)
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
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
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

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat  string
		instanceLabel string
		user          *mongodbflex.User
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
			name: "empty user",
			args: args{
				user: &mongodbflex.User{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.instanceLabel, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
