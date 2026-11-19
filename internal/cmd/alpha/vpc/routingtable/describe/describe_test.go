package describe

import (
	"context"
	"strings"
	"testing"
	"time"

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
		VpcId:          testVpcId,
		RoutingTableId: testRoutingTableId,
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
		valid      bool
		expected   *inputModel
	}{
		{
			name:       "base",
			args:       []string{testRoutingTableId},
			flagValues: fixtureFlagValues(),
			valid:      true,
			expected:   fixtureInputModel(),
		},
		{
			name: "missing project",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, globalflags.ProjectIdFlag)
			}),
		},
		{
			name: "invalid project",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[globalflags.ProjectIdFlag] = "invalid"
			}),
		},
		{
			name: "missing region",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, globalflags.RegionFlag)
			}),
		},
		{
			name: "missing vpc",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, vpcIdFlag)
			}),
		},
		{
			name: "invalid vpc",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[vpcIdFlag] = "invalid"
			}),
		},
		{
			name: "empty vpc",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[vpcIdFlag] = ""
			}),
		},
		{
			name:       "missing routing-table id",
			args:       []string{},
			flagValues: fixtureFlagValues(),
		},
		{
			name:       "invalid routing-table id",
			args:       []string{"invalid"},
			flagValues: fixtureFlagValues(),
		},
		{
			name:       "empty routing-table id",
			args:       []string{""},
			flagValues: fixtureFlagValues(),
		},
		{
			name:       "too many arguments",
			args:       []string{testRoutingTableId, testRoutingTableId},
			flagValues: fixtureFlagValues(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expected, tt.args, tt.flagValues, tt.valid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	ctx := context.Background()
	client := &iaas.DefaultAPIServiceMock{}

	for _, tt := range []struct {
		name  string
		model *inputModel
		want  iaas.ApiGetVPCRoutingTableRequest
	}{
		{
			name:  "default",
			model: fixtureInputModel(),
			want:  client.GetVPCRoutingTable(ctx, testProjectId, testVpcId, "eu02", testRoutingTableId),
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
	createdAt := time.Date(2026, time.October, 1, 12, 0, 0, 0, time.UTC)
	resp := &iaas.VPCRoutingTable{
		Id:          &testRoutingTableId,
		Name:        "rt",
		Description: new("my routing table"),
		Labels: map[string]any{
			"team": "network",
			"env":  "dev",
		},
		DynamicRoutes: new(false),
		SystemRoutes:  new(true),
		CreatedAt:     &createdAt,
		UpdatedAt:     &createdAt,
	}

	for _, tt := range []struct {
		name    string
		format  string
		resp    *iaas.VPCRoutingTable
		want    []string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "empty",
			resp: &iaas.VPCRoutingTable{},
			want: []string{"ID", "NAME", "SYSTEM ROUTES", "DYNAMIC ROUTES"},
		},
		{
			name: "table",
			resp: resp,
			want: []string{testRoutingTableId, "rt", "my routing table", "env: dev", "team: network", "true", "false", "2026-10-01"},
		},
		{
			name:   "json",
			format: print.JSONOutputFormat,
			resp:   resp,
			want:   []string{`"name": "rt"`, `"dynamicRoutes": false`, `"systemRoutes": true`, testRoutingTableId},
		},
		{
			name:   "yaml",
			format: print.YAMLOutputFormat,
			resp:   resp,
			want:   []string{"name: rt", "dynamicRoutes: false", "systemRoutes: true", testRoutingTableId},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.format, tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}

			output := params.Out.String()
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Errorf("output %q does not contain %q", output, want)
				}
			}

			if tt.format == "" && tt.resp != nil && len(tt.resp.Labels) > 0 && strings.Index(output, "env: dev") > strings.Index(output, "team: network") {
				t.Errorf("labels are not sorted: %s", output)
			}
		})
	}
}
