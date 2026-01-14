// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package describe

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
	commonValidation "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/validation"
	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testProjectId = uuid.NewString()
	testRegion    = "eu01"

	testInstanceId  = "instance"
	testDisplayName = "test"
)

// mockExecutable is a mock for the Executable interface
type mockExecutable struct {
	executeFails    bool
	executeNotFound bool
	executeResp     *edge.Instance
}

func (m *mockExecutable) Execute() (*edge.Instance, error) {
	if m.executeFails {
		return nil, errors.New("API error")
	}
	if m.executeNotFound {
		return nil, &oapierror.GenericOpenAPIError{
			StatusCode: http.StatusNotFound,
		}
	}
	return m.executeResp, nil
}

// mockAPIClient is a mock for the edge.APIClient interface
type mockAPIClient struct {
	getInstanceMock       edge.ApiGetInstanceRequest
	getInstanceByNameMock edge.ApiGetInstanceByNameRequest
}

func (m *mockAPIClient) GetInstance(_ context.Context, _, _, _ string) edge.ApiGetInstanceRequest {
	if m.getInstanceMock != nil {
		return m.getInstanceMock
	}
	return &mockExecutable{}
}

func (m *mockAPIClient) GetInstanceByName(_ context.Context, _, _, _ string) edge.ApiGetInstanceByNameRequest {
	if m.getInstanceByNameMock != nil {
		return m.getInstanceByNameMock
	}
	return &mockExecutable{}
}

// Unused methods to satisfy the interface
func (m *mockAPIClient) PostInstances(_ context.Context, _, _ string) edge.ApiPostInstancesRequest {
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
func (m *mockAPIClient) GetTokenByInstanceId(_ context.Context, _, _, _ string) edge.ApiGetTokenByInstanceIdRequest {
	return nil
}
func (m *mockAPIClient) GetTokenByInstanceName(_ context.Context, _, _, _ string) edge.ApiGetTokenByInstanceNameRequest {
	return nil
}

func (m *mockAPIClient) ListPlansProject(_ context.Context, _ string) edge.ApiListPlansProjectRequest {
	return nil
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

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
			name: "instanceId missing",
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
			name:    "instanceId empty",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.InstanceIdFlag] = ""
				}),
			},
		},
		{
			name:    "instanceId too long",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.InstanceIdFlag] = "invalid-instance-id"
				}),
			},
		},
		{
			name:    "instanceId too short",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.InstanceIdFlag] = "id"
				}),
			},
		},
		{
			name:    "name too short",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.InstanceIdFlag)
					flagValues[commonInstance.DisplayNameFlag] = "foo"
				}),
			},
		},
		{
			name:    "name too long",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.InstanceIdFlag)
					flagValues[commonInstance.DisplayNameFlag] = "foofoofoo"
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
	}

	tests := []struct {
		name    string
		wantErr error
		want    *edge.Instance
		args    args
	}{
		{
			name: "get by id success",
			want: &edge.Instance{
				Id:          &testInstanceId,
				DisplayName: &testDisplayName,
			},
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					getInstanceMock: &mockExecutable{
						executeResp: &edge.Instance{
							Id:          &testInstanceId,
							DisplayName: &testDisplayName,
						},
					},
				},
			},
		},
		{
			name: "get by name success",
			want: &edge.Instance{
				Id:          &testInstanceId,
				DisplayName: &testDisplayName,
			},
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					getInstanceByNameMock: &mockExecutable{
						executeResp: &edge.Instance{
							Id:          &testInstanceId,
							DisplayName: &testDisplayName,
						},
					},
				},
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
			name:    "instance not found error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					getInstanceMock: &mockExecutable{
						executeNotFound: true,
					},
				},
			},
		},
		{
			name:    "get by id API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					getInstanceMock: &mockExecutable{
						executeFails: true,
					},
				},
			},
		},
		{
			name:    "get by name API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					getInstanceByNameMock: &mockExecutable{
						executeFails: true,
					},
				},
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
			got, err := run(testCtx, tt.args.model, tt.args.client)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}

			testUtils.AssertValue(t, got, tt.want)
		})
	}
}

func TestOutputResult(t *testing.T) {
	type outputArgs struct {
		model    *inputModel
		instance *edge.Instance
	}

	tests := []struct {
		name    string
		wantErr error
		args    outputArgs
	}{
		{
			name:    "no instance",
			wantErr: &commonErr.NoInstanceError{},
			args: outputArgs{
				model:    fixtureByIdInputModel(),
				instance: nil,
			},
		},
		{
			name: "output json",
			args: outputArgs{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
					model.identifier = nil
				}),
				instance: &edge.Instance{},
			},
		},
		{
			name: "output yaml",
			args: outputArgs{
				model: fixtureInputModel(false, func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
					model.identifier = nil
				}),
				instance: &edge.Instance{},
			},
		},
		{
			name: "output default",
			args: outputArgs{
				model:    fixtureByIdInputModel(),
				instance: &edge.Instance{Id: &testInstanceId},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&types.CmdParams{Printer: p})

			err := outputResult(p, tt.args.model.OutputFormat, tt.args.instance)
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
		want    *describeRequestSpec
		args    args
	}{
		{
			name: "get by id",
			want: &describeRequestSpec{
				ProjectID:  testProjectId,
				Region:     testRegion,
				InstanceId: testInstanceId,
			},
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					getInstanceMock: &mockExecutable{},
				},
			},
		},
		{
			name: "get by name",
			want: &describeRequestSpec{
				ProjectID:    testProjectId,
				Region:       testRegion,
				InstanceName: testDisplayName,
			},
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					getInstanceByNameMock: &mockExecutable{},
				},
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
			testUtils.AssertValue(t, got, tt.want, testUtils.WithIgnoreFields(describeRequestSpec{}, "Execute"))
		})
	}
}
