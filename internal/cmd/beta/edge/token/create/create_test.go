package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testClient     = &edge.APIClient{DefaultAPI: &edge.DefaultAPIService{}}
)

const (
	testRegion     = "eu01"
	testExpiration = 3600
	tokenString    = "test-token-string"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceIdFlag:            testInstanceId,
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
		Expiration: uint64(testExpiration), // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
		InstanceId: testInstanceId,
	}

	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *edge.ApiGetTokenByInstanceIdRequest)) edge.ApiGetTokenByInstanceIdRequest {
	request := testClient.DefaultAPI.GetTokenByInstanceId(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.ExpirationSeconds(int64(testExpiration))
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
			description: "with expiration",
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Expiration = uint64(3600)
			}),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expirationFlag] = "1h"
			}),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no flag values",
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
			description: "invalid expiration format",
			isValid:     false,
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expirationFlag] = "invalid"
			}),
		},
		{
			description: "expiration too short",
			isValid:     false,
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expirationFlag] = "1s"
			}),
		},
		{
			description: "expiration too long",
			isValid:     false,
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expirationFlag] = "13M"
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, func(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
				return parseInput(p, cmd)
			}, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		expectedRequest edge.ApiGetTokenByInstanceIdRequest
		model           *inputModel
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "expiration time",
			model: fixtureInputModel(func(model *inputModel) {
				model.Expiration = uint64(2592000)
			}),
			expectedRequest: fixtureRequest().ExpirationSeconds(int64(2592000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				t.Fatalf("cannot create request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, edge.DefaultAPIService{}),
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
		model *inputModel
		token *edge.Token
	}
	tests := []struct {
		description string
		wantErr     any
		args        args
	}{
		{
			description: "default output format",
			args: args{
				model: fixtureInputModel(),
				token: &edge.Token{
					Token: tokenString,
				},
			},
		},
		{
			description: "JSON output format",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				token: &edge.Token{
					Token: tokenString,
				},
			},
		},
		{
			description: "YAML output format",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
				}),
				token: &edge.Token{
					Token: tokenString,
				},
			},
		},
		{
			description: "nil token",
			wantErr:     true,
			args: args{
				model: fixtureInputModel(),
				token: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.args.model.OutputFormat, tt.args.token)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}
