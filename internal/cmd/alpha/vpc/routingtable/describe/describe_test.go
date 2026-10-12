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

func TestParseInput(t *testing.T) {
	for _, tt := range []struct {
		name   string
		flag   string
		value  string
		remove bool
		args   []string
		valid  bool
	}{
		{
			name:  "base",
			valid: true,
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
			name: "missing routing-table id",
			args: []string{},
		},
		{
			name: "invalid routing-table id",
			args: []string{"invalid"},
		},
		{
			name: "empty routing-table id",
			args: []string{""},
		},
		{
			name: "too many arguments",
			args: []string{testRoutingTableId, testRoutingTableId},
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
			inputArgs := tt.args
			if inputArgs == nil {
				inputArgs = []string{testRoutingTableId}
			}

			var expected *inputModel
			if tt.valid {
				expected = &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{
						ProjectId: testProjectId,
						Region:    "eu02",
						Verbosity: globalflags.VerbosityDefault,
					},
					VpcId:          testVpcId,
					RoutingTableId: testRoutingTableId,
				}
			}

			testutils.TestParseInput(t, NewCmd, parseInput, expected, inputArgs, flagValues, tt.valid)
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
			name: "default",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    "eu02",
				},
				VpcId:          testVpcId,
				RoutingTableId: testRoutingTableId,
			},
			want: client.GetVPCRoutingTable(ctx, testProjectId, testVpcId, "eu02", testRoutingTableId),
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
