package update

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
			VpcId:          testVpcId,
			RoutingTableId: testRoutingTableId,
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
			name:  "name only",
			flag:  nameFlag,
			value: "rt",
			want:  defaultModel(func(m *inputModel) { m.Name = new("rt") }),
		},
		{
			name:  "description only",
			flag:  descriptionFlag,
			value: "my routing table",
			want:  defaultModel(func(m *inputModel) { m.Description = new("my routing table") }),
		},
		{
			name: "clear description",
			flag: descriptionFlag,
			want: defaultModel(func(m *inputModel) { m.Description = new("") }),
		},
		{
			name:  "labels only",
			flag:  labelFlag,
			value: "env=dev,team=network",
			want: defaultModel(func(m *inputModel) {
				m.Labels = map[string]any{
					"env":  "dev",
					"team": "network",
				}
			}),
		},
		{
			name: "empty labels",
			flag: labelFlag,
		},
		{
			name:  "disable dynamic routes",
			flag:  dynamicRoutesFlag,
			value: "false",
			want:  defaultModel(func(m *inputModel) { m.DynamicRoutes = new(false) }),
		},
		{
			name:  "enable dynamic routes",
			flag:  dynamicRoutesFlag,
			value: "true",
			want:  defaultModel(func(m *inputModel) { m.DynamicRoutes = new(true) }),
		},
		{
			name:  "disable system routes",
			flag:  systemRoutesFlag,
			value: "false",
			want:  defaultModel(func(m *inputModel) { m.SystemRoutes = new(false) }),
		},
		{
			name:  "enable system routes",
			flag:  systemRoutesFlag,
			value: "true",
			want:  defaultModel(func(m *inputModel) { m.SystemRoutes = new(true) }),
		},
		{
			name:   "no update flags",
			flag:   nameFlag,
			remove: true,
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
			name:  "invalid labels",
			flag:  labelFlag,
			value: "invalid",
		},
		{
			name:  "invalid boolean",
			flag:  dynamicRoutesFlag,
			value: "invalid",
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
			if tt.want == nil {
				flagValues[nameFlag] = "rt"
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

			testutils.TestParseInput(t, NewCmd, parseInput, tt.want, inputArgs, flagValues, tt.want != nil)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	defaultModel := func(opts ...func(*inputModel)) *inputModel {
		model := &inputModel{
			GlobalFlagModel: &globalflags.GlobalFlagModel{
				ProjectId: testProjectId,
				Region:    "eu02",
			},
			VpcId:          testVpcId,
			RoutingTableId: testRoutingTableId,
			Name:           new("rt"),
			Description:    new("description"),
			Labels:         map[string]any{"env": "dev"},
			DynamicRoutes:  new(true),
			SystemRoutes:   new(true),
		}
		for _, opt := range opts {
			opt(model)
		}
		return model
	}
	defaultPayload := func(opts ...func(*iaas.UpdateVPCRoutingTablePayload)) iaas.UpdateVPCRoutingTablePayload {
		payload := iaas.UpdateVPCRoutingTablePayload{
			Name:          new("rt"),
			Description:   new("description"),
			Labels:        map[string]any{"env": "dev"},
			DynamicRoutes: new(true),
			SystemRoutes:  new(true),
		}
		for _, opt := range opts {
			opt(&payload)
		}
		return payload
	}

	for _, tt := range []struct {
		name        string
		model       *inputModel
		wantPayload iaas.UpdateVPCRoutingTablePayload
	}{
		{
			name:        "all fields",
			model:       defaultModel(),
			wantPayload: defaultPayload(),
		},
		{
			name:        "empty labels",
			model:       defaultModel(func(m *inputModel) { m.Labels = map[string]any{} }),
			wantPayload: defaultPayload(func(p *iaas.UpdateVPCRoutingTablePayload) { p.Labels = map[string]any{} }),
		},
		{
			name:        "omitted labels",
			model:       defaultModel(func(m *inputModel) { m.Labels = nil }),
			wantPayload: defaultPayload(func(p *iaas.UpdateVPCRoutingTablePayload) { p.Labels = nil }),
		},
		{
			name: "name only leaves other fields unchanged",
			model: defaultModel(func(m *inputModel) {
				m.Labels = nil
				m.Description = nil
				m.DynamicRoutes = nil
				m.SystemRoutes = nil
			}),
			wantPayload: iaas.UpdateVPCRoutingTablePayload{
				Name: new("rt"),
			},
		},
		{
			name: "disabled routes",
			model: defaultModel(func(m *inputModel) {
				m.DynamicRoutes = new(false)
				m.SystemRoutes = new(false)
			}),
			wantPayload: defaultPayload(func(p *iaas.UpdateVPCRoutingTablePayload) {
				p.DynamicRoutes = new(false)
				p.SystemRoutes = new(false)
			}),
		},
		{
			name:        "clear description",
			model:       defaultModel(func(m *inputModel) { m.Description = new("") }),
			wantPayload: defaultPayload(func(p *iaas.UpdateVPCRoutingTablePayload) { p.Description = new("") }),
		},
		{
			name:        "omitted name",
			model:       defaultModel(func(m *inputModel) { m.Name = nil }),
			wantPayload: defaultPayload(func(p *iaas.UpdateVPCRoutingTablePayload) { p.Name = nil }),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client := &iaas.DefaultAPIServiceMock{}

			want := client.UpdateVPCRoutingTable(ctx, testProjectId, testVpcId, "eu02", testRoutingTableId).UpdateVPCRoutingTablePayload(tt.wantPayload)
			got := buildRequest(ctx, tt.model, client)
			if diff := cmp.Diff(want, got, cmp.AllowUnexported(want), cmpopts.EquateComparable(ctx, iaas.DefaultAPIService{})); diff != "" {
				t.Fatalf("request mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	for _, tt := range []struct {
		name    string
		format  string
		resp    *iaas.VPCRoutingTable
		want    string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name:    "missing id",
			resp:    &iaas.VPCRoutingTable{},
			wantErr: true,
		},
		{
			name: "empty id",
			resp: &iaas.VPCRoutingTable{
				Id: new(""),
			},
			wantErr: true,
		},
		{
			name: "updated",
			resp: &iaas.VPCRoutingTable{
				Id: &testRoutingTableId,
			},
			want: "Updated routing-table with ID \"" + testRoutingTableId + "\".",
		},
		{
			name:   "json",
			format: print.JSONOutputFormat,
			resp: &iaas.VPCRoutingTable{
				Id:   &testRoutingTableId,
				Name: "rt",
			},
			want: `"name": "rt"`,
		},
		{
			name:   "yaml",
			format: print.YAMLOutputFormat,
			resp: &iaas.VPCRoutingTable{
				Id:   &testRoutingTableId,
				Name: "rt",
			},
			want: "name: rt",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.format, tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !strings.Contains(params.Out.String(), tt.want) {
				t.Errorf("output %q does not contain %q", params.Out.String(), tt.want)
			}
		})
	}
}
