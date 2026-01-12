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
	executeResp  *edge.InstanceList
}

func (m *mockExecutable) Execute() (*edge.InstanceList, error) {
	if m.executeFails {
		return nil, errors.New("API error")
	}

	if m.executeResp != nil {
		return m.executeResp, nil
	}
	return &edge.InstanceList{
		Instances: &[]edge.Instance{
			{Id: utils.Ptr("instance-1"), DisplayName: utils.Ptr("namea")},
			{Id: utils.Ptr("instance-2"), DisplayName: utils.Ptr("nameb")},
		},
	}, nil
}

// mockAPIClient is a mock for the edge.APIClient interface
type mockAPIClient struct {
	getInstancesMock edge.ApiGetInstancesRequest
}

func (m *mockAPIClient) GetInstances(_ context.Context, _, _ string) edge.ApiGetInstancesRequest {
	if m.getInstancesMock != nil {
		return m.getInstancesMock
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
			name: "success",
			want: fixtureInputModel(),
			args: args{
				flags: fixtureFlagValues(),
				cmpOpts: []testUtils.ValueComparisonOption{
					testUtils.WithAllowUnexported(inputModel{}),
				},
			},
		},
		{
			name: "with limit",
			want: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(10))
			}),
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[limitFlag] = "10"
				}),
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
			name:    "limit invalid",
			wantErr: "invalid syntax",
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[limitFlag] = "invalid"
				}),
			},
		},
		{
			name:    "limit less than 1",
			wantErr: &cliErr.FlagValidationError{},
			args: args{
				flags: fixtureFlagValues(func(flagValues map[string]string) {
					flagValues[limitFlag] = "0"
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
		want    []edge.Instance
		args    args
	}{
		{
			name: "list success",
			want: []edge.Instance{
				{Id: utils.Ptr("instance-1"), DisplayName: utils.Ptr("namea")},
				{Id: utils.Ptr("instance-2"), DisplayName: utils.Ptr("nameb")},
			},
			args: args{
				model:  fixtureInputModel(),
				client: &mockAPIClient{},
			},
		},
		{
			name: "list success with limit",
			want: []edge.Instance{
				{Id: utils.Ptr("instance-1"), DisplayName: utils.Ptr("namea")},
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
			want: []edge.Instance{
				{Id: utils.Ptr("instance-1"), DisplayName: utils.Ptr("namea")},
				{Id: utils.Ptr("instance-2"), DisplayName: utils.Ptr("nameb")},
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
			want: []edge.Instance{},
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					getInstancesMock: &mockExecutable{
						executeResp: &edge.InstanceList{Instances: &[]edge.Instance{}},
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
					getInstancesMock: &mockExecutable{
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
		instances    []edge.Instance
		projectLabel string
	}

	tests := []struct {
		name    string
		wantErr error
		args    args
	}{
		{
			name: "no instance",
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
				instances: []edge.Instance{
					{Id: utils.Ptr("instance-1"), DisplayName: utils.Ptr("namea")},
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
				instances: []edge.Instance{
					{Id: utils.Ptr("instance-1"), DisplayName: utils.Ptr("namea")},
				},
				projectLabel: "test-project",
			},
		},
		{
			name: "output default with instances",
			args: args{
				model: fixtureInputModel(),
				instances: []edge.Instance{
					{
						Id:          utils.Ptr("instance-1"),
						DisplayName: utils.Ptr("namea"),
						FrontendUrl: utils.Ptr("https://example.com"),
					},
					{
						Id:          utils.Ptr("instance-2"),
						DisplayName: utils.Ptr("nameb"),
						FrontendUrl: utils.Ptr("https://example2.com"),
					},
				},
				projectLabel: "test-project",
			},
		},
		{
			name: "output default with no instances",
			args: args{
				model:        fixtureInputModel(),
				instances:    []edge.Instance{},
				projectLabel: "test-project",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&types.CmdParams{Printer: p})

			err := outputResult(p, tt.args.model.OutputFormat, tt.args.projectLabel, tt.args.instances)
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
				Region:    testRegion,
			},
			args: args{
				model: fixtureInputModel(),
				client: &mockAPIClient{
					getInstancesMock: &mockExecutable{},
				},
			},
		},
		{
			name: "success with limit",
			want: &listRequestSpec{
				ProjectID: testProjectId,
				Region:    testRegion,
				Limit:     utils.Ptr(int64(10)),
			},
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Limit = utils.Ptr(int64(10))
				}),
				client: &mockAPIClient{
					getInstancesMock: &mockExecutable{},
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
