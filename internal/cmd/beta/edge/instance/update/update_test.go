// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package update

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
	testCtx         = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testProjectId   = uuid.NewString()
	testRegion      = "eu01"
	testInstanceId  = "instance"
	testDisplayName = "test"
	testDescription = "new description"
	testPlanId      = uuid.NewString()
)

type mockExecutable struct {
	executeFails                bool
	executeNotFound             bool
	capturedUpdatePayload       *edge.UpdateInstancePayload
	capturedUpdateByNamePayload *edge.UpdateInstanceByNamePayload
}

func (m *mockExecutable) Execute() error {
	if m.executeFails {
		return errors.New("API error")
	}
	if m.executeNotFound {
		return &oapierror.GenericOpenAPIError{
			StatusCode: http.StatusNotFound,
		}
	}
	return nil
}

func (m *mockExecutable) UpdateInstancePayload(payload edge.UpdateInstancePayload) edge.ApiUpdateInstanceRequest {
	if m.capturedUpdatePayload != nil {
		*m.capturedUpdatePayload = payload
	}
	return m
}

func (m *mockExecutable) UpdateInstanceByNamePayload(payload edge.UpdateInstanceByNamePayload) edge.ApiUpdateInstanceByNameRequest {
	if m.capturedUpdateByNamePayload != nil {
		*m.capturedUpdateByNamePayload = payload
	}
	return m
}

type mockAPIClient struct {
	updateInstanceMock       edge.ApiUpdateInstanceRequest
	updateInstanceByNameMock edge.ApiUpdateInstanceByNameRequest
}

func (m *mockAPIClient) UpdateInstance(_ context.Context, _, _, _ string) edge.ApiUpdateInstanceRequest {
	if m.updateInstanceMock != nil {
		return m.updateInstanceMock
	}
	return &mockExecutable{}
}

func (m *mockAPIClient) UpdateInstanceByName(_ context.Context, _, _, _ string) edge.ApiUpdateInstanceByNameRequest {
	if m.updateInstanceByNameMock != nil {
		return m.updateInstanceByNameMock
	}
	return &mockExecutable{}
}

// Unused methods to satisfy the interface
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
		globalflags.ProjectIdFlag:      testProjectId,
		globalflags.RegionFlag:         testRegion,
		commonInstance.InstanceIdFlag:  testInstanceId,
		commonInstance.DescriptionFlag: testDescription,
		commonInstance.PlanIdFlag:      testPlanId,
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
		Description: &testDescription,
		PlanId:      &testPlanId,
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
			name:    "no update flags",
			wantErr: true,
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.DescriptionFlag)
					delete(flagValues, commonInstance.PlanIdFlag)
				}),
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
			name:    "plan id invalid",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.PlanIdFlag] = "not-a-uuid"
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

func TestBuildRequest(t *testing.T) {
	type args struct {
		model  *inputModel
		client client.APIClient
	}
	tests := []struct {
		name    string
		args    args
		want    *updateRequestSpec
		wantErr error
	}{
		{
			name: "by id",
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					updateInstanceMock: &mockExecutable{},
				},
			},
			want: &updateRequestSpec{
				ProjectID:  testProjectId,
				Region:     testRegion,
				InstanceId: testInstanceId,
				Payload: edge.UpdateInstancePayload{
					Description: &testDescription,
					PlanId:      &testPlanId,
				},
			},
		},
		{
			name: "by name",
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					updateInstanceByNameMock: &mockExecutable{},
				},
			},
			want: &updateRequestSpec{
				ProjectID:    testProjectId,
				Region:       testRegion,
				InstanceName: testDisplayName,
				Payload: edge.UpdateInstancePayload{
					Description: &testDescription,
					PlanId:      &testPlanId,
				},
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
				testUtils.AssertValue(t, got, tt.want, testUtils.WithIgnoreFields(updateRequestSpec{}, "Execute"))
			}
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
		args    args
	}{
		{
			name: "update by id success",
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					updateInstanceMock: &mockExecutable{},
				},
			},
		},
		{
			name: "update by name success",
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					updateInstanceByNameMock: &mockExecutable{},
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
					updateInstanceMock: &mockExecutable{
						executeNotFound: true,
					},
				},
			},
		},
		{
			name:    "update by id API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model: fixtureByIdInputModel(),
				client: &mockAPIClient{
					updateInstanceMock: &mockExecutable{
						executeFails: true,
					},
				},
			},
		},
		{
			name:    "update by name API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model: fixtureByNameInputModel(),
				client: &mockAPIClient{
					updateInstanceByNameMock: &mockExecutable{
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
			err := run(testCtx, tt.args.model, tt.args.client)
			testUtils.AssertError(t, err, tt.wantErr)
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
