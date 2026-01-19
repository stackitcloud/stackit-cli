// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package create

import (
	"context"
	"errors"
	"net/http"
	"testing"

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

// mockTokenWaiter is a mock for the tokenWaiter interface
type mockTokenWaiter struct {
	waitFails    bool
	waitNotFound bool
	waitResp     *edge.Token
}

func (m *mockTokenWaiter) WaitWithContext(_ context.Context) (*edge.Token, error) {
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

	// Default token response
	tokenString := "test-token-string"
	return &edge.Token{
		Token: &tokenString,
	}, nil
}

// testWaiterFactoryProvider is a test implementation that returns mock waiters.
type testWaiterFactoryProvider struct {
	waiter tokenWaiter
}

func (t *testWaiterFactoryProvider) getTokenWaiter(_ context.Context, model *inputModel, _ client.APIClient) (tokenWaiter, error) {
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
func (m *mockAPIClient) GetTokenByInstanceId(_ context.Context, _, _, _ string) edge.ApiGetTokenByInstanceIdRequest {
	return nil
}

func (m *mockAPIClient) GetTokenByInstanceName(_ context.Context, _, _, _ string) edge.ApiGetTokenByInstanceNameRequest {
	return nil
}

func (m *mockAPIClient) ListPlansProject(_ context.Context, _ string) edge.ApiListPlansProjectRequest {
	return nil
}

func (m *mockAPIClient) CreateInstance(_ context.Context, _, _ string) edge.ApiCreateInstanceRequest {
	return nil
}
func (m *mockAPIClient) GetInstance(_ context.Context, _, _, _ string) edge.ApiGetInstanceRequest {
	return nil
}
func (m *mockAPIClient) GetInstanceByName(_ context.Context, _, _, _ string) edge.ApiGetInstanceByNameRequest {
	return nil
}
func (m *mockAPIClient) ListInstances(_ context.Context, _, _ string) edge.ApiListInstancesRequest {
	return nil
}
func (m *mockAPIClient) UpdateInstance(_ context.Context, _, _, _ string) edge.ApiUpdateInstanceRequest {
	return nil
}
func (m *mockAPIClient) UpdateInstanceByName(_ context.Context, _, _, _ string) edge.ApiUpdateInstanceByNameRequest {
	return nil
}
func (m *mockAPIClient) DeleteInstance(_ context.Context, _, _, _ string) edge.ApiDeleteInstanceRequest {
	return nil
}
func (m *mockAPIClient) DeleteInstanceByName(_ context.Context, _, _, _ string) edge.ApiDeleteInstanceByNameRequest {
	return nil
}
func (m *mockAPIClient) GetKubeconfigByInstanceId(_ context.Context, _, _, _ string) edge.ApiGetKubeconfigByInstanceIdRequest {
	return nil
}
func (m *mockAPIClient) GetKubeconfigByInstanceName(_ context.Context, _, _, _ string) edge.ApiGetKubeconfigByInstanceNameRequest {
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
		Expiration: uint64(commonKubeconfig.ExpirationSecondsDefault), // Default 1 hour
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
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.InstanceIdFlag] = ""
				}),
			},
		},
		{
			name:    "instance id too long",
			wantErr: &cliErr.FlagValidationError{},
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
		waiter tokenWaiter
	}
	tests := []struct {
		name      string
		wantErr   any
		wantToken bool
		args      args
	}{
		{
			name:      "run by id success",
			wantToken: true,
			args: args{
				model:  fixtureByIdInputModel(),
				client: &mockAPIClient{},
				waiter: &mockTokenWaiter{},
			},
		},
		{
			name:      "run by name success",
			wantToken: true,
			args: args{
				model:  fixtureByNameInputModel(),
				client: &mockAPIClient{},
				waiter: &mockTokenWaiter{},
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
				waiter: &mockTokenWaiter{},
			},
		},
		{
			name:    "instance not found error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model:  fixtureByIdInputModel(),
				client: &mockAPIClient{},
				waiter: &mockTokenWaiter{waitNotFound: true},
			},
		},
		{
			name:    "get token by id API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model:  fixtureByIdInputModel(),
				client: &mockAPIClient{},
				waiter: &mockTokenWaiter{waitFails: true},
			},
		},
		{
			name:    "get token by name API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model:  fixtureByNameInputModel(),
				client: &mockAPIClient{},
				waiter: &mockTokenWaiter{waitFails: true},
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
				waiter: &mockTokenWaiter{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override production waiterProvider package level variable for testing
			prodWaiterProvider := waiterProvider
			waiterProvider = &testWaiterFactoryProvider{waiter: tt.args.waiter}
			defer func() { waiterProvider = prodWaiterProvider }()

			got, err := run(testCtx, tt.args.model, tt.args.client)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}
			if tt.wantToken && got == nil {
				t.Fatal("expected non-nil token")
			}
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
		want    bool
		wantErr error
		args    args
	}{
		{
			name: "by id",
			want: true,
			args: args{model: fixtureByIdInputModel()},
		},
		{
			name: "by name",
			want: true,
			args: args{model: fixtureByNameInputModel()},
		},
		{
			name:    "no id or name",
			wantErr: &commonErr.NoIdentifierError{},
			args: args{model: fixtureInputModel(false, func(model *inputModel) {
				model.identifier = nil
			})},
		},
		{
			name:    "unknown identifier",
			wantErr: &commonErr.InvalidIdentifierError{},
			args: args{model: fixtureInputModel(false, func(model *inputModel) {
				model.identifier.Flag = "unknown"
			})},
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
		model *inputModel
		token *edge.Token
	}
	tests := []struct {
		name    string
		wantErr any
		args    args
	}{
		{
			name: "default output format",
			args: args{
				model: fixtureByIdInputModel(),
				token: &edge.Token{
					Token: func() *string { s := "test-token"; return &s }(),
				},
			},
		},
		{
			name: "JSON output format",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				token: &edge.Token{
					Token: func() *string { s := "test-token"; return &s }(),
				},
			},
		},
		{
			name: "YAML output format",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
				}),
				token: &edge.Token{
					Token: func() *string { s := "test-token"; return &s }(),
				},
			},
		},
		{
			name:    "nil token",
			wantErr: true,
			args: args{
				model: fixtureByIdInputModel(),
				token: nil,
			},
		},
		{
			name:    "nil token string",
			wantErr: true,
			args: args{
				model: fixtureByIdInputModel(),
				token: &edge.Token{Token: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&types.CmdParams{Printer: p})
			err := outputResult(p, tt.args.model.OutputFormat, tt.args.token)
			testUtils.AssertError(t, err, tt.wantErr)
		})
	}
}
