package list

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

var (
	testProjectId      = uuid.NewString()
	testVpcId          = uuid.NewString()
	testRoutingTableId = uuid.NewString()
)

func TestParseInput(t *testing.T) {
	defaultModel := func(opts ...func(*inputModel)) *inputModel {
		model := &inputModel{
			GlobalFlagModel: &globalflags.GlobalFlagModel{
				ProjectId: testProjectId,
				Region:    "eu02",
				Verbosity: globalflags.VerbosityDefault,
			},
			VpcId: testVpcId,
		}
		for _, opt := range opts {
			opt(model)
		}
		return model
	}

	for _, tt := range []struct {
		name   string
		flag   string
		value  string
		remove bool
		args   []string
		want   *inputModel
	}{
		{
			name: "defaults",
			want: defaultModel(),
		},
		{
			name:  "label selector",
			flag:  labelSelectorFlag,
			value: "env=dev",
			want:  defaultModel(func(m *inputModel) { m.LabelSelector = new("env=dev") }),
		},
		{
			name:  "filter",
			flag:  filterFlag,
			value: `name == "rt"`,
			want:  defaultModel(func(m *inputModel) { m.Filter = new(`name == "rt"`) }),
		},
		{
			name:  "limit",
			flag:  limitFlag,
			value: "10",
			want:  defaultModel(func(m *inputModel) { m.Limit = new(int64(10)) }),
		},
		{
			name:  "zero limit",
			flag:  limitFlag,
			value: "0",
		},
		{
			name:  "negative limit",
			flag:  limitFlag,
			value: "-1",
		},
		{
			name:  "invalid limit",
			flag:  limitFlag,
			value: "invalid",
		},
		{
			name:   "missing project",
			flag:   globalflags.ProjectIdFlag,
			remove: true,
		},
		{
			name:  "invalid project",
			flag:  globalflags.ProjectIdFlag,
			value: "invalid",
		},
		{
			name:   "missing region",
			flag:   globalflags.RegionFlag,
			remove: true,
		},
		{
			name:   "missing vpc",
			flag:   vpcIdFlag,
			remove: true,
		},
		{
			name:  "invalid vpc",
			flag:  vpcIdFlag,
			value: "invalid",
		},
		{
			name: "empty vpc",
			flag: vpcIdFlag,
		},
		{
			name: "unexpected argument",
			args: []string{"unexpected"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			flagValues := map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    "eu02",
				vpcIdFlag:                 testVpcId,
			}
			if tt.flag != "" {
				flagValues[tt.flag] = tt.value
				if tt.remove {
					delete(flagValues, tt.flag)
				}
			}

			testutils.TestParseInput(t, NewCmd, parseInput, tt.want, tt.args, flagValues, tt.want != nil)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	ctx := context.Background()
	client := &iaas.DefaultAPIServiceMock{}
	defaultModel := func(opts ...func(*inputModel)) *inputModel {
		model := &inputModel{
			GlobalFlagModel: &globalflags.GlobalFlagModel{
				ProjectId: testProjectId,
				Region:    "eu02",
			},
			VpcId: testVpcId,
		}
		for _, opt := range opts {
			opt(model)
		}
		return model
	}

	for _, tt := range []struct {
		name  string
		model *inputModel
		want  iaas.ApiListVPCRoutingTablesRequest
	}{
		{
			name:  "default",
			model: defaultModel(),
			want:  client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02"),
		},
		{
			name:  "label selector",
			model: defaultModel(func(m *inputModel) { m.LabelSelector = new("env=dev") }),
			want:  client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").LabelSelector("env=dev"),
		},
		{
			name:  "filter",
			model: defaultModel(func(m *inputModel) { m.Filter = new(`name == "rt"`) }),
			want:  client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").Filter(`name == "rt"`),
		},
		{
			name: "combined filters",
			model: defaultModel(func(m *inputModel) {
				m.LabelSelector = new("env=dev")
				m.Filter = new(`name == "rt"`)
			}),
			want: client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").LabelSelector("env=dev").Filter(`name == "rt"`),
		},
		{
			name: "empty filters",
			model: defaultModel(func(m *inputModel) {
				m.LabelSelector = new("")
				m.Filter = new("")
			}),
			want: client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").LabelSelector("").Filter(""),
		},
		{
			name:  "limit is client side",
			model: defaultModel(func(m *inputModel) { m.Limit = new(int64(10)) }),
			want:  client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRequest(ctx, tt.model, client)
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(tt.want), cmpopts.EquateComparable(ctx, iaas.DefaultAPIService{})); diff != "" {
				t.Fatalf("request mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	resp := &iaas.VPCRoutingTableList{
		Items: []iaas.VPCRoutingTable{
			{
				Id:          &testRoutingTableId,
				Name:        "rt",
				Description: new("my routing table"),
				Labels: map[string]any{
					"team": "network",
					"env":  "dev",
				},
				SystemRoutes:  new(true),
				DynamicRoutes: new(false),
			},
		},
	}

	for _, tt := range []struct {
		name    string
		format  string
		resp    *iaas.VPCRoutingTableList
		want    []string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "nil items",
			resp: &iaas.VPCRoutingTableList{},
			want: []string{"No routing-tables found", testVpcId},
		},
		{
			name: "empty items",
			resp: &iaas.VPCRoutingTableList{
				Items: []iaas.VPCRoutingTable{},
			},
			want: []string{"No routing-tables found", testVpcId},
		},
		{
			name: "table",
			resp: resp,
			want: []string{testRoutingTableId, "rt", "my routing table", "env: dev", "team: network", "true", "false"},
		},
		{
			name: "missing optional fields",
			resp: &iaas.VPCRoutingTableList{
				Items: []iaas.VPCRoutingTable{
					{
						Name: "rt",
					},
				},
			},
			want: []string{"rt", "SYSTEM ROUTES", "DYNAMIC ROUTES"},
		},
		{
			name:   "json",
			format: print.JSONOutputFormat,
			resp:   resp,
			want:   []string{`"items"`, `"name": "rt"`, `"dynamicRoutes": false`, testRoutingTableId},
		},
		{
			name:   "yaml",
			format: print.YAMLOutputFormat,
			resp:   resp,
			want:   []string{"items:", "name: rt", "dynamicRoutes: false", testRoutingTableId},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.format, testVpcId, tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			output := params.Out.String()
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Errorf("output %q does not contain %q", output, want)
				}
			}
			if tt.format == "" && strings.Contains(output, "env: dev") && strings.Index(output, "env: dev") > strings.Index(output, "team: network") {
				t.Errorf("labels are not sorted: %s", output)
			}
		})
	}
}
