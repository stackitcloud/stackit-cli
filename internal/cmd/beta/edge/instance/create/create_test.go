// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package create

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	commonErr "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/error"
	commonInstance "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/instance"
	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testProjectId = uuid.NewString()
	testRegion    = "eu01"

	testName        = "test"
	testPlanId      = uuid.NewString()
	testDescription = "Initial instance description"
	testInstanceId  = uuid.NewString()
)

// mockExecutable is a mock for the Executable interface used by the SDK
type mockExecutable struct {
	executeFails bool
	resp         *edge.Instance
}

func (m *mockExecutable) PostInstancesPayload(_ edge.PostInstancesPayload) edge.ApiPostInstancesRequest {
	// This method is needed to satisfy the interface. It allows chaining in buildRequest.
	return m
}
func (m *mockExecutable) Execute() (*edge.Instance, error) {
	if m.executeFails {
		return nil, errors.New("API error")
	}
	if m.resp != nil {
		return m.resp, nil
	}
	return &edge.Instance{Id: &testInstanceId}, nil
}

// mockAPIClient is a mock for the client.APIClient interface
type mockAPIClient struct {
	postInstancesMock edge.ApiPostInstancesRequest
}

func (m *mockAPIClient) PostInstances(_ context.Context, _, _ string) edge.ApiPostInstancesRequest {
	if m.postInstancesMock != nil {
		return m.postInstancesMock
	}
	return &mockExecutable{}
}

// Unused methods to satisfy the client.APIClient interface
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
		commonInstance.DisplayNameFlag: testName,
		commonInstance.DescriptionFlag: testDescription,
		commonInstance.PlanIdFlag:      testPlanId,
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
		DisplayName: testName,
		Description: testDescription,
		PlanId:      testPlanId,
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
			name: "create success",
			want: fixtureInputModel(),
			args: args{
				flags: fixtureFlagValues(),
				cmpOpts: []testUtils.ValueComparisonOption{
					testUtils.WithAllowUnexported(inputModel{}, globalflags.GlobalFlagModel{}),
				},
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
			name:    "name missing",
			wantErr: "required flag(s) \"name\" not set",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					delete(flagValues, commonInstance.DisplayNameFlag)
				}),
			},
		},
		{
			name:    "name too long",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.DisplayNameFlag] = "this-name-is-way-too-long-for-the-validation"
				}),
			},
		},
		{
			name:    "name too short",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.DisplayNameFlag] = "in"
				}),
			},
		},
		{
			name:    "name invalid",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.DisplayNameFlag] = "1test"
				}),
			},
		},
		{
			name:    "plan invalid",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.PlanIdFlag] = "invalid-uuid"
				}),
			},
		},
		{
			name:    "description too long",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[commonInstance.DescriptionFlag] = strings.Repeat("a", 257)
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
		name string
		args args
		want *createRequestSpec
	}{
		{
			name: "success",
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					postInstancesMock: &mockExecutable{},
				},
			},
			want: &createRequestSpec{
				ProjectID: testProjectId,
				Region:    testRegion,
				Payload: edge.PostInstancesPayload{
					DisplayName: &testName,
					Description: &testDescription,
					PlanId:      &testPlanId,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := buildRequest(testCtx, tt.args.model, tt.args.client)

			if got != nil {
				if got.Execute == nil {
					t.Error("expected non-nil Execute function")
				}
				testUtils.AssertValue(t, got, tt.want, testUtils.WithIgnoreFields(createRequestSpec{}, "Execute"))
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
		want    *edge.Instance
		args    args
	}{
		{
			name: "create success",
			want: &edge.Instance{Id: &testInstanceId},
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					postInstancesMock: &mockExecutable{
						resp: &edge.Instance{Id: &testInstanceId},
					},
				},
			},
		},
		{
			name:    "create API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					postInstancesMock: &mockExecutable{
						executeFails: true,
					},
				},
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
	type args struct {
		model        *inputModel
		instance     *edge.Instance
		projectLabel string
	}

	tests := []struct {
		name    string
		wantErr error
		args    args
	}{
		{
			name:    "no instance",
			wantErr: &commonErr.NoInstanceError{},
			args: args{
				model: fixtureInputModel(),
			},
		},
		{
			name: "output json",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				instance: &edge.Instance{},
			},
		},
		{
			name: "output yaml",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
				}),
				instance: &edge.Instance{},
			},
		},
		{
			name: "output default",
			args: args{
				model:    fixtureInputModel(),
				instance: &edge.Instance{Id: &testInstanceId},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&types.CmdParams{Printer: p})

			err := outputResult(p, tt.args.model.OutputFormat, tt.args.model.Async, tt.args.projectLabel, tt.args.instance)
			testUtils.AssertError(t, err, tt.wantErr)
		})
	}
}
