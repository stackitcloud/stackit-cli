package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &alb.APIClient{}

var (
	testProjectId     = uuid.NewString()
	testRegion        = "eu01"
	testCredentialRef = "credential-12345"
	testDisplayname   = "displayname"
	testUsername      = "testuser"
	testPassword      = "testpassword"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		usernameFlag:              testUsername,
		displaynameFlag:           testDisplayname,
		passwordFlag:              testPassword,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) inputModel {
	model := inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
			Region:    testRegion,
		},
		Username:       &testUsername,
		Displayname:    &testDisplayname,
		CredentialsRef: &testCredentialRef,
		Password:       &testPassword,
	}
	for _, mod := range mods {
		mod(&model)
	}
	return model
}

func fixtureRequest(mods ...func(request *alb.ApiUpdateCredentialsRequest)) alb.ApiUpdateCredentialsRequest {
	request := testClient.UpdateCredentials(testCtx, testProjectId, testRegion, testCredentialRef)
	request = request.UpdateCredentialsPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *alb.UpdateCredentialsPayload)) alb.UpdateCredentialsPayload {
	payload := alb.UpdateCredentialsPayload{
		DisplayName: &testDisplayname,
		Password:    &testPassword,
		Username:    &testUsername,
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
		args          []string
		isValid       bool
		expectedModel inputModel
	}{
		{
			description:   "base",
			args:          []string{testCredentialRef},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			args:        []string{testCredentialRef},
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
			},
			isValid: false,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Username = nil
				model.Displayname = nil
			}),
		},
		{
			description: "required values",
			args:        []string{testCredentialRef},
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				usernameFlag:              testUsername,
				displaynameFlag:           testDisplayname,
				passwordFlag:              testPassword,
			},
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&types.CmdParams{Printer: p})
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err = cmd.Flags().Set(flag, value)
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

			model := parseInput(p, cmd, tt.args)

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
		model           inputModel
		expectedRequest alb.ApiUpdateCredentialsRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, &tt.model, testClient)
			if err != nil {
				t.Fatalf("cannot build request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmp.AllowUnexported(alb.NullableString{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		item  *alb.UpdateCredentialsResponse
		model inputModel
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				item:  nil,
				model: fixtureInputModel(),
			},
			wantErr: true,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
