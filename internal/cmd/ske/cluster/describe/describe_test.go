package describe

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
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{}
var testProjectId = uuid.NewString()
var testClusterName = "cluster"

const testRegion = "eu01"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testClusterName,
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
		ClusterName: testClusterName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *ske.ApiGetClusterRequest)) ske.ApiGetClusterRequest {
	request := testClient.GetCluster(testCtx, testProjectId, testRegion, testClusterName)
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
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
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
		expectedRequest ske.ApiGetClusterRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
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
		outputFormat string
		cluster      *ske.Cluster
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
			name: "empty cluster",
			args: args{
				cluster: &ske.Cluster{},
			},
			wantErr: false,
		},
		{
			name: "cluster with single error",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{
							{
								Code:    utils.Ptr("SKE_INFRA_SNA_NETWORK_NOT_FOUND"),
								Message: utils.Ptr("Network configuration not found"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cluster with multiple errors",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{
							{
								Code:    utils.Ptr("SKE_INFRA_SNA_NETWORK_NOT_FOUND"),
								Message: utils.Ptr("Network configuration not found"),
							},
							{
								Code:    utils.Ptr("SKE_NODE_MACHINE_TYPE_NOT_FOUND"),
								Message: utils.Ptr("Specified machine type unavailable"),
							},
							{
								Code:    utils.Ptr("SKE_FETCHING_ERRORS_NOT_POSSIBLE"),
								Message: utils.Ptr("Fetching errors not possible"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cluster with error but no message",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{
							{
								Code: utils.Ptr("SKE_FETCHING_ERRORS_NOT_POSSIBLE"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cluster with nil errors",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Status: &ske.ClusterStatus{
						Errors: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cluster with empty errors array",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cluster without status",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
				},
			},
			wantErr: false,
		},
		{
			name: "JSON output format with errors",
			args: args{
				outputFormat: print.JSONOutputFormat,
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{
							{
								Code:    utils.Ptr("SKE_INFRA_SNA_NETWORK_NOT_FOUND"),
								Message: utils.Ptr("Network configuration not found"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "YAML output format with errors",
			args: args{
				outputFormat: print.YAMLOutputFormat,
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{
							{
								Code:    utils.Ptr("SKE_INFRA_SNA_NETWORK_NOT_FOUND"),
								Message: utils.Ptr("Network configuration not found"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cluster with kubernetes info and errors",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Kubernetes: &ske.Kubernetes{
						Version: utils.Ptr("1.28.0"),
					},
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{
							{
								Code:    utils.Ptr("SKE_INFRA_SNA_NETWORK_NOT_FOUND"),
								Message: utils.Ptr("Network configuration not found"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cluster with extensions and errors",
			args: args{
				outputFormat: "",
				cluster: &ske.Cluster{
					Name: utils.Ptr("test-cluster"),
					Extensions: &ske.Extension{
						Acl: &ske.ACL{
							AllowedCidrs: &[]string{"10.0.0.0/8"},
							Enabled:      utils.Ptr(true),
						},
					},
					Status: &ske.ClusterStatus{
						Errors: &[]ske.ClusterError{
							{
								Code:    utils.Ptr("SKE_INFRA_SNA_NETWORK_NOT_FOUND"),
								Message: utils.Ptr("Network configuration not found"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.cluster); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
