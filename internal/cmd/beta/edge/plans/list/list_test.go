// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package list

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testProjectId = uuid.NewString()
	testRegion    = "eu01"
)

// mockExecutable is a mock for the Executable interface
type mockExecutable struct {
	executeFails bool
	executeResp  *edge.PlanList
}

func (m *mockExecutable) Execute() (*edge.PlanList, error) {
	if m.executeFails {
		return nil, errors.New("API error")
	}

	if m.executeResp != nil {
		return m.executeResp, nil
	}
	return &edge.PlanList{
		ValidPlans: &[]edge.Plan{
			{Id: utils.Ptr("plan-1"), Name: utils.Ptr("Standard")},
			{Id: utils.Ptr("plan-2"), Name: utils.Ptr("Premium")},
		},
	}, nil
}

// mockAPIClient is a mock for the edge.APIClient interface
type mockAPIClient struct {
	getPlansMock edge.ApiListPlansProjectRequest
}

func (m *mockAPIClient) ListPlansProject(_ context.Context, _ string) edge.ApiListPlansProjectRequest {
	if m.getPlansMock != nil {
		return m.getPlansMock
	}
	return &mockExecutable{}
}

// Unused methods to satisfy the interface
func (m *mockAPIClient) CreateInstance(_ context.Context, _, _ string) edge.ApiCreateInstanceRequest {
	return nil
}
func (m *mockAPIClient) GetInstance(_ context.Context, _, _, _ string) edge.ApiGetInstanceRequest {
	return nil
}

func (m *mockAPIClient) ListInstances(_ context.Context, _, _ string) edge.ApiListInstancesRequest {
	return nil
}

func (m *mockAPIClient) GetInstanceByName(_ context.Context, _, _, _ string) edge.ApiGetInstanceByNameRequest {
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

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		limitFlag:                 "10",
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
		Limit: utils.Ptr(int64(10)),
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
			name: "list success",
			want: fixtureInputModel(),
			args: args{
				flags: fixtureFlagValues(),
				cmpOpts: []testUtils.ValueComparisonOption{
					testUtils.WithAllowUnexported(inputModel{}),
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
			name:    "limit invalid value",
			wantErr: "invalid syntax",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[limitFlag] = "invalid"
				}),
			},
		},
		{
			name:    "limit is zero",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[limitFlag] = "0"
				}),
			},
		},
		{
			name:    "limit is negative",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[limitFlag] = "-0"
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
		want    []edge.Plan
		args    args
	}{
		{
			name: "list success",
			want: []edge.Plan{
				{Id: utils.Ptr("plan-1"), Name: utils.Ptr("Standard")},
				{Id: utils.Ptr("plan-2"), Name: utils.Ptr("Premium")},
			},
			args: args{
				model:  fixtureInputModel(),
				client: &mockAPIClient{},
			},
		},
		{
			name: "list success with limit",
			want: []edge.Plan{
				{Id: utils.Ptr("plan-1"), Name: utils.Ptr("Standard")},
			},
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Limit = utils.Ptr(int64(1))
				}),
				client: &mockAPIClient{},
			},
		},
		{
			name: "list success with limit greater than items",
			want: []edge.Plan{
				{Id: utils.Ptr("plan-1"), Name: utils.Ptr("Standard")},
				{Id: utils.Ptr("plan-2"), Name: utils.Ptr("Premium")},
			},
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Limit = utils.Ptr(int64(5))
				}),
				client: &mockAPIClient{},
			},
		},
		{
			name: "list success with no items",
			want: []edge.Plan{},
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					getPlansMock: &mockExecutable{
						executeResp: &edge.PlanList{ValidPlans: &[]edge.Plan{}},
					},
				},
			},
		},
		{
			name:    "list API error",
			wantErr: &cliErr.RequestFailedError{},
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					getPlansMock: &mockExecutable{
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
		plans        []edge.Plan
		projectLabel string
	}

	tests := []struct {
		name    string
		wantErr error
		args    args
	}{
		{
			name: "output json",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				plans: []edge.Plan{
					{Id: utils.Ptr("plan-1"), Name: utils.Ptr("Standard")},
				},
				projectLabel: "test-project",
			},
		},
		{
			name: "output yaml",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
				}),
				plans: []edge.Plan{
					{Id: utils.Ptr("plan-1"), Name: utils.Ptr("Standard")},
				},
				projectLabel: "test-project",
			},
		},
		{
			name: "output default with plans",
			args: args{
				model: fixtureInputModel(),
				plans: []edge.Plan{
					{
						Id:          utils.Ptr("plan-1"),
						Name:        utils.Ptr("Standard"),
						Description: utils.Ptr("Standard plan description"),
					},
					{
						Id:          utils.Ptr("plan-2"),
						Name:        utils.Ptr("Premium"),
						Description: utils.Ptr("Premium plan description"),
					},
				},
				projectLabel: "test-project",
			},
		},
		{
			name: "output default with no plans",
			args: args{
				model:        fixtureInputModel(),
				plans:        []edge.Plan{},
				projectLabel: "test-project",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&types.CmdParams{Printer: p})

			err := outputResult(p, tt.args.model.OutputFormat, tt.args.projectLabel, tt.args.plans)
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
		want    *listRequestSpec
		args    args
	}{
		{
			name: "success",
			want: &listRequestSpec{
				ProjectID: testProjectId,
			},
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Limit = nil
				}),
				client: &mockAPIClient{
					getPlansMock: &mockExecutable{},
				},
			},
		},
		{
			name: "success with limit",
			want: &listRequestSpec{
				ProjectID: testProjectId,
				Limit:     utils.Ptr(int64(10)),
			},
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					getPlansMock: &mockExecutable{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildRequest(testCtx, tt.args.model, tt.args.client)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}
			testUtils.AssertValue(t, got, tt.want, testUtils.WithIgnoreFields(listRequestSpec{}, "Execute"))
		})
	}
}
