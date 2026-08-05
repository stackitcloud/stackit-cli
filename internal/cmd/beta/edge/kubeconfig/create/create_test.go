package create

import (
	"context"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
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
	testKubeconfig = `
apiVersion: v1
clusters:
- cluster:
    server: https://server-1.com
  name: cluster-1
contexts:
- context:
    cluster: cluster-1
    user: user-1
  name: context-1
current-context: context-1
kind: Config
preferences: {}
users:
- name: user-1
  user: {}
`
)

// Helper function to create a new instance of Kubeconfig
//
//nolint:gocritic // ptrToRefParam: Required by edge.Kubeconfig API which expects *map[string]interface{}
func testKubeconfigMap() map[string]interface{} {
	var kubeconfigMap map[string]interface{}
	err := yaml.Unmarshal([]byte(testKubeconfig), &kubeconfigMap)
	if err != nil {
		// This should never happen in tests with valid YAML
		panic(err)
	}
	return kubeconfigMap
}

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
		DisableWriting: false,
		Filepath:       nil,
		Overwrite:      false,
		Expiration:     uint64(3600), // Default 1 hour
		SwitchContext:  false,
		InstanceId:     testInstanceId,
	}

	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *edge.ApiGetKubeconfigByInstanceIdRequest)) edge.ApiGetKubeconfigByInstanceIdRequest {
	request := testClient.DefaultAPI.GetKubeconfigByInstanceId(testCtx, testProjectId, testRegion, testInstanceId)
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
			description: "disable writing and invalid output format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[disableWritingFlag] = "true"
			}),
			isValid: false,
		},
		{
			description: "disable writing and valid output format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[disableWritingFlag] = "true"
				flagValues[globalflags.OutputFormatFlag.Name()] = print.YAMLOutputFormat
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DisableWriting = true
				model.OutputFormat = print.YAMLOutputFormat
			}),
			isValid: true,
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
		{
			description: "enable overwrite",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[overwriteFlag] = "true"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Overwrite = true
			}),
			isValid: true,
		},
		{
			description: "disable overwrite",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[overwriteFlag] = "false"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Overwrite = false
			}),
			isValid: true,
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
		expectedRequest edge.ApiGetKubeconfigByInstanceIdRequest
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
		model      *inputModel
		kubeconfig *edge.Kubeconfig
	}

	tests := []struct {
		name    string
		wantErr any
		args    args
	}{
		{
			name:    "no kubeconfig",
			wantErr: true,
			args: args{
				model:      fixtureInputModel(),
				kubeconfig: nil,
			},
		},
		{
			name:    "kubeconfig with nil kubeconfig data",
			wantErr: true,
			args: args{
				model:      fixtureInputModel(),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: nil},
			},
		},
		{
			name: "output json with disable writing",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "output yaml with disable writing",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "output default with disable writing",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "file writing enabled (default behavior)",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.AssumeYes = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "file writing with overwrite enabled",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Overwrite = true
					model.AssumeYes = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "file writing with switch context enabled",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.SwitchContext = true
					model.AssumeYes = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.args.model.OutputFormat, tt.args.model, tt.args.kubeconfig)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}
