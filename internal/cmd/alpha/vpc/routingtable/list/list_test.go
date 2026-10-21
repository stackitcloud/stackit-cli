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

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    "eu02",
		vpcIdFlag:                 testVpcId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}

	return flagValues
}

func fixtureInputModel(mods ...func(*inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    "eu02",
			Verbosity: globalflags.VerbosityDefault,
		},
		VpcId: testVpcId,
	}
	for _, mod := range mods {
		mod(model)
	}

	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		flagValues map[string]string
		want       *inputModel
	}{
		{
			name:       "defaults",
			flagValues: fixtureFlagValues(),
			want:       fixtureInputModel(),
		},
		{
			name: "label selector",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[labelSelectorFlag] = "env=dev"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.LabelSelector = new("env=dev")
			}),
		},
		{
			name: "filter",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[filterFlag] = `name == "rt"`
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Filter = new(`name == "rt"`)
			}),
		},
		{
			name: "limit",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[limitFlag] = "10"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Limit = new(int64(10))
			}),
		},
		{
			name: "combined options",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[labelSelectorFlag] = "env=dev"
				f[filterFlag] = `name == "rt"`
				f[limitFlag] = "10"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.LabelSelector = new("env=dev")
				m.Filter = new(`name == "rt"`)
				m.Limit = new(int64(10))
			}),
		},
		{
			name: "zero limit",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[limitFlag] = "0"
			}),
		},
		{
			name: "negative limit",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[limitFlag] = "-1"
			}),
		},
		{
			name: "invalid limit",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[limitFlag] = "invalid"
			}),
		},
		{
			name: "missing project",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, globalflags.ProjectIdFlag)
			}),
		},
		{
			name: "invalid project",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[globalflags.ProjectIdFlag] = "invalid"
			}),
		},
		{
			name: "missing region",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, globalflags.RegionFlag)
			}),
		},
		{
			name: "missing vpc",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, vpcIdFlag)
			}),
		},
		{
			name: "invalid vpc",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[vpcIdFlag] = "invalid"
			}),
		},
		{
			name: "empty vpc",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[vpcIdFlag] = ""
			}),
		},
		{
			name:       "unexpected argument",
			args:       []string{"unexpected"},
			flagValues: fixtureFlagValues(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.want, tt.args, tt.flagValues, tt.want != nil)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	ctx := context.Background()
	client := &iaas.DefaultAPIServiceMock{}

	for _, tt := range []struct {
		name  string
		model *inputModel
		want  iaas.ApiListVPCRoutingTablesRequest
	}{
		{
			name:  "default",
			model: fixtureInputModel(),
			want:  client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02"),
		},
		{
			name: "label selector",
			model: fixtureInputModel(func(m *inputModel) {
				m.LabelSelector = new("env=dev")
			}),
			want: client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").LabelSelector("env=dev"),
		},
		{
			name: "filter",
			model: fixtureInputModel(func(m *inputModel) {
				m.Filter = new(`name == "rt"`)
			}),
			want: client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").Filter(`name == "rt"`),
		},
		{
			name: "combined filters",
			model: fixtureInputModel(func(m *inputModel) {
				m.LabelSelector = new("env=dev")
				m.Filter = new(`name == "rt"`)
			}),
			want: client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").LabelSelector("env=dev").Filter(`name == "rt"`),
		},
		{
			name: "empty filters",
			model: fixtureInputModel(func(m *inputModel) {
				m.LabelSelector = new("")
				m.Filter = new("")
			}),
			want: client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02").LabelSelector("").Filter(""),
		},
		{
			name: "limit is client side",
			model: fixtureInputModel(func(m *inputModel) {
				m.Limit = new(int64(10))
			}),
			want: client.ListVPCRoutingTables(ctx, testProjectId, testVpcId, "eu02"),
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
