// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package delete

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

// mockExecutable implements the SDK delete request interface for testing.
type mockExecutable struct {
	executeFails    bool
	executeNotFound bool
}

func (m *mockExecutable) Execute() error {
	if m.executeNotFound {
		return &oapierror.GenericOpenAPIError{
			StatusCode: http.StatusNotFound,
			Body:       []byte(`{"message":"not found"}`),
		}
	}
	if m.executeFails {
		return errors.New("execute failed")
	}
	return nil
}

// mockAPIClient provides the minimal API client behavior required by the tests.
type mockAPIClient struct {
	deleteInstanceMock       edge.ApiDeleteInstanceRequest
	deleteInstanceByNameMock edge.ApiDeleteInstanceByNameRequest
}

func (m *mockAPIClient) DeleteInstance(_ context.Context, _, _, _ string) edge.ApiDeleteInstanceRequest {
	if m.deleteInstanceMock != nil {
		return m.deleteInstanceMock
	}
	return &mockExecutable{}
}

func (m *mockAPIClient) DeleteInstanceByName(_ context.Context, _, _, _ string) edge.ApiDeleteInstanceByNameRequest {
	if m.deleteInstanceByNameMock != nil {
		return m.deleteInstanceByNameMock
	}
	return &mockExecutable{}
}

// Unused methods to satisfy the client.APIClient interface.
func (m *mockAPIClient) PostInstances(_ context.Context, _, _ string) edge.ApiPostInstancesRequest {
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

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
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

func fixtureInputModel(useDisplayName bool, mods ...func(*inputModel)) *inputModel {
	identifier := &commonValidation.Identifier{
		Flag:  commonInstance.InstanceIdFlag,
		Value: testInstanceId,
	}
	if useDisplayName {
		identifier = &commonValidation.Identifier{
			Flag:  commonInstance.DisplayNameFlag,
			Value: testDisplayName,
		}
	}

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		identifier: identifier,
	}

	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureByIdInputModel(mods ...func(*inputModel)) *inputModel {
	return fixtureInputModel(false, mods...)
}

func fixtureByNameInputModel(mods ...func(*inputModel)) *inputModel {
	return fixtureInputModel(true, mods...)
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
					testUtils.WithAllowUnexported(inputModel{}, globalflags.GlobalFlagModel{}),
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
					testUtils.WithAllowUnexported(inputModel{}, globalflags.GlobalFlagModel{}),
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
		args    args
		wantErr error
	}{
		{
			name: "delete by id success",
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					deleteInstanceMock: &mockExecutable{},
				},
			},
		},
		{
			name: "delete by id API error",
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					deleteInstanceMock: &mockExecutable{executeFails: true},
				},
			},
			wantErr: &cliErr.RequestFailedError{},
		},
		{
			name: "delete by id not found",
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					deleteInstanceMock: &mockExecutable{executeNotFound: true},
				},
			},
			wantErr: &cliErr.RequestFailedError{},
		},
		{
			name: "delete by name success",
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					deleteInstanceByNameMock: &mockExecutable{},
				},
			},
		},
		{
			name: "delete by name API error",
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					deleteInstanceByNameMock: &mockExecutable{executeFails: true},
				},
			},
			wantErr: &cliErr.RequestFailedError{},
		},
		{
			name: "delete by name not found",
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					deleteInstanceByNameMock: &mockExecutable{executeNotFound: true},
				},
			},
			wantErr: &cliErr.RequestFailedError{},
		},
		{
			name: "no identifier",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.identifier = nil
				}),
				client: &mockAPIClient{},
			},
			wantErr: &commonErr.NoIdentifierError{},
		},
		{
			name: "invalid identifier",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.identifier = &commonValidation.Identifier{Flag: "unknown", Value: "value"}
				}),
				client: &mockAPIClient{},
			},
			wantErr: &cliErr.BuildRequestError{},
		},
		{
			name: "nil model",
			args: args{
				model:  nil,
				client: &mockAPIClient{},
			},
			wantErr: &commonErr.NoIdentifierError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(testCtx, tt.args.model, tt.args.client)
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
		args    args
		want    *deleteRequestSpec
		wantErr error
	}{
		{
			name: "by id",
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					deleteInstanceMock: &mockExecutable{},
				},
			},
			want: &deleteRequestSpec{
				ProjectID:  testProjectId,
				Region:     testRegion,
				InstanceId: testInstanceId,
			},
		},
		{
			name: "by name",
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					deleteInstanceByNameMock: &mockExecutable{},
				},
			},
			want: &deleteRequestSpec{
				ProjectID:    testProjectId,
				Region:       testRegion,
				InstanceName: testDisplayName,
			},
		},
		{
			name: "no identifier",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.identifier = nil
				}),
				client: &mockAPIClient{},
			},
			wantErr: &commonErr.NoIdentifierError{},
		},
		{
			name: "invalid identifier",
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.identifier = &commonValidation.Identifier{Flag: "unknown", Value: "val"}
				}),
				client: &mockAPIClient{},
			},
			wantErr: &cliErr.BuildRequestError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildRequest(testCtx, tt.args.model, tt.args.client)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}
			if got != nil {
				if got.Execute == nil {
					t.Error("expected non-nil Execute function")
				}
				testUtils.AssertValue(t, got, tt.want, testUtils.WithIgnoreFields(deleteRequestSpec{}, "Execute"))
			}
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
			name: "by id identifier",
			want: true,
			args: args{
				model: fixtureByIdInputModel(),
			},
		},
		{
			name: "by name identifier",
			want: true,
			args: args{
				model: fixtureByNameInputModel(),
			},
		},
		{
			name:    "nil model",
			wantErr: &commonErr.NoIdentifierError{},
			want:    false,
			args: args{
				model: nil,
			},
		},
		{
			name:    "nil identifier",
			wantErr: &commonErr.NoIdentifierError{},
			want:    false,
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.identifier = nil
				}),
			},
		},
		{
			name:    "invalid identifier",
			wantErr: &commonErr.InvalidIdentifierError{},
			want:    false,
			args: args{
				model: fixtureByIdInputModel(func(model *inputModel) {
					model.identifier = &commonValidation.Identifier{Flag: "unsupported", Value: "value"}
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
