package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"
)

const (
	testRegion = "eu01"

	testDisplayName = "display-name"
	testDescription = "description"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &logs.APIClient{}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		instanceIdFlag:  testInstanceId,
		displayNameFlag: testDisplayName,
		descriptionFlag: testDescription,
		permissionsFlag: "read,write",
		lifetimeFlag:    "0",
	}
	for _, mod := range mods {
		mod(flagValues)
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

		InstanceId:  testInstanceId,
		Description: utils.Ptr(testDescription),
		DisplayName: testDisplayName,
		Lifetime:    utils.Ptr(int64(0)),
		Permissions: []string{
			"read",
			"write",
		},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *logs.ApiCreateAccessTokenRequest)) logs.ApiCreateAccessTokenRequest {
	request := testClient.CreateAccessToken(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.CreateAccessTokenPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *logs.CreateAccessTokenPayload)) logs.CreateAccessTokenPayload {
	payload := logs.CreateAccessTokenPayload{
		DisplayName: utils.Ptr(testDisplayName),
		Description: utils.Ptr(testDescription),
		Lifetime:    utils.Ptr(int64(0)),
		Permissions: utils.Ptr([]string{
			"read",
			"write",
		}),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, lifetimeFlag)
				delete(flagValues, descriptionFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Lifetime = nil
				model.Description = nil
			}),
		},
		{
			description: "one permission",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[permissionsFlag] = "read"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Permissions = []string{
					"read",
				}
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
			description: "instance id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "lifetime invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[lifetimeFlag] = "invalid-integer"
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
	var tests = []struct {
		description     string
		model           *inputModel
		expectedRequest logs.ApiCreateAccessTokenRequest
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
	type args struct {
		outputFormat  string
		instanceLabel string
		accessToken   *logs.AccessToken
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				instanceLabel: "",
				accessToken: utils.Ptr(logs.AccessToken{
					Id: utils.Ptr(uuid.NewString()),
					Permissions: utils.Ptr([]string{
						"read",
						"write",
					}),
					DisplayName: utils.Ptr("Token"),
					AccessToken: utils.Ptr("Secret access token"),
					Creator:     utils.Ptr(uuid.NewString()),
					Expires:     utils.Ptr(false),
					Status:      utils.Ptr(logs.ACCESSTOKENSTATUS_ACTIVE),
				}),
			},
			wantErr: false,
		},
		{
			name: "empty access token",
			args: args{
				instanceLabel: "",
				accessToken:   utils.Ptr(logs.AccessToken{}),
			},
			wantErr: false,
		},
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.instanceLabel, tt.args.accessToken); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
