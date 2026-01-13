// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package create

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	commonErr "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/error"
	commonInstance "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/instance"
	commonKubeconfig "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/kubeconfig"
	commonValidation "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/validation"
	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
)

type testCtxKey struct{}

var (
	testCtx         = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testProjectId   = uuid.NewString()
	testRegion      = "eu01"
	testInstanceId  = "instance"
	testDisplayName = "test"
	testExpiration  = "1h"
)

const (
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
func testKubeconfigMap() *map[string]interface{} {
	var kubeconfigMap map[string]interface{}
	err := yaml.Unmarshal([]byte(testKubeconfig), &kubeconfigMap)
	if err != nil {
		// This should never happen in tests with valid YAML
		panic(err)
	}
	return utils.Ptr(kubeconfigMap)
}

// mockKubeconfigWaiter is a mock for the kubeconfigWaiter interface
type mockKubeconfigWaiter struct {
	waitFails    bool
	waitNotFound bool
	waitResp     *edge.Kubeconfig
}

func (m *mockKubeconfigWaiter) WaitWithContext(_ context.Context) (*edge.Kubeconfig, error) {
	if m.waitFails {
		return nil, errors.New("wait error")
	}
	if m.waitNotFound {
		return nil, &oapierror.GenericOpenAPIError{
			StatusCode: http.StatusNotFound,
		}
	}
	if m.waitResp != nil {
		return m.waitResp, nil
	}

	// Default kubeconfig response
	return &edge.Kubeconfig{
		Kubeconfig: testKubeconfigMap(),
	}, nil
}

// testWaiterFactoryProvider is a test implementation that returns mock waiters.
type testWaiterFactoryProvider struct {
	waiter kubeconfigWaiter
}

func (t *testWaiterFactoryProvider) getKubeconfigWaiter(_ context.Context, model *inputModel, _ client.APIClient) (kubeconfigWaiter, error) {
	if model == nil || model.identifier == nil {
		return nil, &commonErr.NoIdentifierError{}
	}

	// Validate identifier like the real implementation
	switch model.identifier.Flag {
	case commonInstance.InstanceIdFlag, commonInstance.DisplayNameFlag:
		// Return our mock waiter directly, bypassing the client type casting issue
		return t.waiter, nil
	default:
		return nil, commonErr.NewInvalidIdentifierError(model.identifier.Flag)
	}
}

// mockAPIClient is a mock for the edge.APIClient interface
type mockAPIClient struct{}

// Unused methods to satisfy the interface
func (m *mockAPIClient) GetKubeconfigByInstanceId(_ context.Context, _, _, _ string) edge.ApiGetKubeconfigByInstanceIdRequest {
	return nil
}

func (m *mockAPIClient) GetKubeconfigByInstanceName(_ context.Context, _, _, _ string) edge.ApiGetKubeconfigByInstanceNameRequest {
	return nil
}

func (m *mockAPIClient) GetTokenByInstanceId(_ context.Context, _, _, _ string) edge.ApiGetTokenByInstanceIdRequest {
	return nil
}

func (m *mockAPIClient) GetTokenByInstanceName(_ context.Context, _, _, _ string) edge.ApiGetTokenByInstanceNameRequest {
	return nil
}

func (m *mockAPIClient) ListPlansProject(_ context.Context, _ string) edge.ApiListPlansProjectRequest {
	return nil
}

func (m *mockAPIClient) PostInstances(_ context.Context, _, _ string) edge.ApiPostInstancesRequest {
	return nil
}

func (m *mockAPIClient) DeleteInstance(_ context.Context, _, _, _ string) edge.ApiDeleteInstanceRequest {
	return nil
}

func (m *mockAPIClient) DeleteInstanceByName(_ context.Context, _, _, _ string) edge.ApiDeleteInstanceByNameRequest {
	return nil
}

func (m *mockAPIClient) GetInstance(_ context.Context, _, _, _ string) edge.ApiGetInstanceRequest {
	return nil
}

func (m *mockAPIClient) GetInstanceByName(_ context.Context, _, _, _ string) edge.ApiGetInstanceByNameRequest {
	return nil
}

func (m *mockAPIClient) GetInstances(_ context.Context, _, _ string) edge.ApiGetInstancesRequest {
	return nil
}

func (m *mockAPIClient) UpdateInstance(_ context.Context, _, _, _ string) edge.ApiUpdateInstanceRequest {
	return nil
}

func (m *mockAPIClient) UpdateInstanceByName(_ context.Context, _, _, _ string) edge.ApiUpdateInstanceByNameRequest {
	return nil
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag:     testProjectId,
		globalflags.RegionFlag:        testRegion,
		commonInstance.InstanceIdFlag: testInstanceId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureByIdInputModel(mods ...func(model *inputModel)) *inputModel {
	return fixtureInputModel(false, mods...)
}

func fixtureByNameInputModel(mods ...func(model *inputModel)) *inputModel {
	return fixtureInputModel(true, mods...)
}

func fixtureInputModel(useName bool, mods ...func(model *inputModel)) *inputModel {
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
	}

	if useName {
		model.identifier = &commonValidation.Identifier{
			Flag:  commonInstance.DisplayNameFlag,
			Value: testDisplayName,
		}
	} else {
		model.identifier = &commonValidation.Identifier{
			Flag:  commonInstance.InstanceIdFlag,
			Value: testInstanceId,
		}
	}

	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	type args struct {
		flags   map[string]string
		cmpOpts []testUtils.ValueComparisonOption
	}

	tests := []struct {
		name    string
		wantErr any
		want    *inputModel
		args    args
	}{
		{
			name: "by id",
			want: fixtureByIdInputModel(),
			args: args{
				flags: fixtureFlagValues(),
				cmpOpts: []testUtils.ValueComparisonOption{
					testUtils.WithAllowUnexported(inputModel{}),
				},
			},
		},
		{
			name: "by name",
			want: fixtureByNameInputModel(),
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.InstanceIdFlag)
					flagValues[commonInstance.DisplayNameFlag] = testDisplayName
				}),
				cmpOpts: []testUtils.ValueComparisonOption{
					testUtils.WithAllowUnexported(inputModel{}),
				},
			},
		},
		{
			name: "with expiration",
			want: fixtureByIdInputModel(func(model *inputModel) {
				model.Expiration = uint64(3600)
			}),
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonKubeconfig.ExpirationFlag] = testExpiration
				}),
				cmpOpts: []testUtils.ValueComparisonOption{
					testUtils.WithAllowUnexported(inputModel{}),
				},
			},
		},
		{
			name:    "by id and name",
			wantErr: true,
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.DisplayNameFlag] = testDisplayName
				}),
			},
		},
		{
			name:    "no flag values",
			wantErr: true,
			args: args{
				flags: map[string]string{},
			},
		},
		{
			name:    "project id missing",
			wantErr: &cliErr.ProjectIdError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, globalflags.ProjectIdFlag)
				}),
			},
		},
		{
			name:    "project id empty",
			wantErr: "value cannot be empty",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[globalflags.ProjectIdFlag] = ""
				}),
			},
		},
		{
			name:    "project id invalid",
			wantErr: "invalid UUID length",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
				}),
			},
		},
		{
			name:    "instance id missing",
			wantErr: true,
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.InstanceIdFlag)
				}),
			},
		},
		{
			name:    "instance id empty",
			wantErr: "id may not be empty",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.InstanceIdFlag] = ""
				}),
			},
		},
		{
			name:    "instance id too long",
			wantErr: "id is too long",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.InstanceIdFlag] = "invalid-instance-id"
				}),
			},
		},
		{
			name:    "instance id too short",
			wantErr: "id is too short",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.InstanceIdFlag] = "id"
				}),
			},
		},
		{
			name:    "name too short",
			wantErr: "name is too short",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.InstanceIdFlag)
					flagValues[commonInstance.DisplayNameFlag] = "foo"
				}),
			},
		},
		{
			name:    "name too long",
			wantErr: "name is too long",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.InstanceIdFlag)
					flagValues[commonInstance.DisplayNameFlag] = "foofoofoo"
				}),
			},
		},
		{
			name:    "disable writing and invalid output format",
			wantErr: "valid output formats for this command are",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonKubeconfig.DisableWritingFlag] = "true"
					flagValues[globalflags.OutputFormatFlag] = print.PrettyOutputFormat
				}),
			},
		},
		{
			name:    "disable writing and default output format",
			wantErr: "must be used with --output-format",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonKubeconfig.DisableWritingFlag] = "true"
				}),
			},
		},
		{
			name: "disable writing and valid output format",
			want: fixtureByIdInputModel(func(model *inputModel) {
				model.DisableWriting = true
				model.OutputFormat = print.YAMLOutputFormat
			}),
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonKubeconfig.DisableWritingFlag] = "true"
					flagValues[globalflags.OutputFormatFlag] = print.YAMLOutputFormat
				}),
				cmpOpts: []testUtils.ValueComparisonOption{
					testUtils.WithAllowUnexported(inputModel{}),
				},
			},
		},
		{
			name:    "invalid expiration format",
			wantErr: "invalid time string format",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonKubeconfig.ExpirationFlag] = "invalid"
				}),
			},
		},
		{
			name:    "expiration too short",
			wantErr: "expiration is too small",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonKubeconfig.ExpirationFlag] = "1s"
				}),
			},
		},
		{
			name:    "expiration too long",
			wantErr: "expiration is too large",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonKubeconfig.ExpirationFlag] = "13M"
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caseOpts := []testUtils.ParseInputCaseOption{}
			if len(tt.args.cmpOpts) > 0 {
				caseOpts = append(caseOpts, testUtils.WithParseInputCmpOptions(tt.args.cmpOpts...))
			}

			testUtils.RunParseInputCase(t, testUtils.ParseInputTestCase[*inputModel]{
				Name:       tt.name,
				Flags:      tt.args.flags,
				WantModel:  tt.want,
				WantErr:    tt.wantErr,
				CmdFactory: NewCmd,
				ParseInputFunc: func(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
					return parseInput(p, cmd)
				},
			}, caseOpts...)
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		model  *inputModel
		client client.APIClient
		waiter kubeconfigWaiter
	}

	tests := []struct {
		name    string
		wantErr error
		args    args
	}{
		{
			name: "run by id success",
			args: args{
				model:  fixtureByIdInputModel(),
				client: &mockAPIClient{},
				waiter: &mockKubeconfigWaiter{},
			},
		},
		{
			name: "run by name success",
			args: args{
				model:  fixtureByNameInputModel(),
				client: &mockAPIClient{},
				waiter: &mockKubeconfigWaiter{},
			},
		},
		{
			name:    "no id or name",
			wantErr: &commonErr.NoIdentifierError{},
			args: args{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.identifier = nil
				}),
				client: &mockAPIClient{},
				waiter: &mockKubeconfigWaiter{},
			},
		},
		{
			name:    "instance not found error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model:  fixtureByIdInputModel(),
				client: &mockAPIClient{},
				waiter: &mockKubeconfigWaiter{waitNotFound: true},
			},
		},
		{
			name:    "get kubeconfig by id API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model:  fixtureByIdInputModel(),
				client: &mockAPIClient{},
				waiter: &mockKubeconfigWaiter{waitFails: true},
			},
		},
		{
			name:    "get kubeconfig by name API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model:  fixtureByNameInputModel(),
				client: &mockAPIClient{},
				waiter: &mockKubeconfigWaiter{waitFails: true},
			},
		},
		{
			name:    "identifier invalid",
			wantErr: &commonErr.InvalidIdentifierError{},
			args: args{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.identifier = &commonValidation.Identifier{
						Flag:  "unknown-flag",
						Value: "some-value",
					}
				}),
				client: &mockAPIClient{},
				waiter: &mockKubeconfigWaiter{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override production waiterProvider package level variable for testing
			prodWaiterProvider := waiterProvider
			waiterProvider = &testWaiterFactoryProvider{waiter: tt.args.waiter}
			defer func() { waiterProvider = prodWaiterProvider }()

			_, err := run(testCtx, tt.args.model, tt.args.client)
			testUtils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	type args struct {
		model  *inputModel
		client client.APIClient
	}

	tests := []struct {
		name    string
		wantErr error
		want    *createRequestSpec
		args    args
	}{
		{
			name: "by id",
			want: &createRequestSpec{
				ProjectID:  testProjectId,
				Region:     testRegion,
				InstanceId: testInstanceId,
				Expiration: int64(commonKubeconfig.ExpirationSecondsDefault),
			},
			args: args{
				model:  fixtureByIdInputModel(),
				client: &mockAPIClient{},
			},
		},
		{
			name: "by name",
			want: &createRequestSpec{
				ProjectID:    testProjectId,
				Region:       testRegion,
				InstanceName: testDisplayName,
				Expiration:   int64(commonKubeconfig.ExpirationSecondsDefault),
			},
			args: args{
				model:  fixtureByNameInputModel(),
				client: &mockAPIClient{},
			},
		},
		{
			name:    "no id or name",
			wantErr: &commonErr.NoIdentifierError{},
			args: args{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.identifier = nil
				}),
				client: &mockAPIClient{},
			},
		},
		{
			name:    "identifier invalid",
			wantErr: &commonErr.InvalidIdentifierError{},
			args: args{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.identifier = &commonValidation.Identifier{
						Flag:  "unknown-flag",
						Value: "some-value",
					}
				}),
				client: &mockAPIClient{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildRequest(testCtx, tt.args.model, tt.args.client)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}
			testUtils.AssertValue(t, got, tt.want, testUtils.WithIgnoreFields(createRequestSpec{}, "Execute"))
		})
	}
}

func TestGetWaiterFactory(t *testing.T) {
	type args struct {
		model *inputModel
	}

	tests := []struct {
		name    string
		wantErr error
		want    bool
		args    args
	}{
		{
			name: "by id",
			want: true,
			args: args{
				model: fixtureByIdInputModel(),
			},
		},
		{
			name: "by name",
			want: true,
			args: args{
				model: fixtureByNameInputModel(),
			},
		},
		{
			name:    "no id or name",
			wantErr: &commonErr.NoIdentifierError{},
			want:    false,
			args: args{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.identifier = nil
				}),
			},
		},
		{
			name:    "unknown identifier",
			wantErr: &commonErr.InvalidIdentifierError{},
			want:    false,
			args: args{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.identifier.Flag = "unknown"
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getWaiterFactory(testCtx, tt.args.model)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}

			if tt.want && got == nil {
				t.Fatal("expected non-nil waiter factory")
			}
			if !tt.want && got != nil {
				t.Fatal("expected nil waiter factory")
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
				model:      fixtureByIdInputModel(),
				kubeconfig: nil,
			},
		},
		{
			name:    "kubeconfig with nil kubeconfig data",
			wantErr: true,
			args: args{
				model:      fixtureByIdInputModel(),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: nil},
			},
		},
		{
			name: "output json with disable writing",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "output yaml with disable writing",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "output default with disable writing",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "output by name with json format and disable writing",
			args: args{
				model: fixtureByNameInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "output by name with yaml format and disable writing",
			args: args{
				model: fixtureByNameInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "output by name default with disable writing",
			args: args{
				model: fixtureByNameInputModel(func(model *inputModel) {
					model.DisableWriting = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "file writing enabled (default behavior)",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.AssumeYes = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "file writing with overwrite enabled",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.Overwrite = true
					model.AssumeYes = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
		{
			name: "file writing with switch context enabled",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.SwitchContext = true
					model.AssumeYes = true
				}),
				kubeconfig: &edge.Kubeconfig{Kubeconfig: testKubeconfigMap()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&types.CmdParams{Printer: p})

			err := outputResult(p, tt.args.model.OutputFormat, tt.args.model, tt.args.kubeconfig)
			testUtils.AssertError(t, err, tt.wantErr)
		})
	}
}
