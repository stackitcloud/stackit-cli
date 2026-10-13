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

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    "eu02",
		vpcIdFlag:                 testVpcId,
		nameFlag:                  "rt",
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
		Name:           new("rt"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixturePayload(mods ...func(*iaas.UpdateVPCRoutingTablePayload)) iaas.UpdateVPCRoutingTablePayload {
	payload := iaas.UpdateVPCRoutingTablePayload{
		Name: new("rt"),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		flagValues map[string]string
		want       *inputModel
	}{
		{
			name:       "name only",
			args:       []string{testRoutingTableId},
			flagValues: fixtureFlagValues(),
			want:       fixtureInputModel(),
		},
		{
			name: "description only",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
				f[descriptionFlag] = "my routing table"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.Description = new("my routing table")
			}),
		},
		{
			name: "clear description",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
				f[descriptionFlag] = ""
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.Description = new("")
			}),
		},
		{
			name: "labels only",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
				f[labelFlag] = "env=dev,team=network"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.Labels = map[string]any{
					"env":  "dev",
					"team": "network",
				}
			}),
		},
		{
			name: "empty labels",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[labelFlag] = ""
			}),
		},
		{
			name: "disable dynamic routes",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
				f[dynamicRoutesFlag] = "false"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.DynamicRoutes = new(false)
			}),
		},
		{
			name: "enable dynamic routes",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
				f[dynamicRoutesFlag] = "true"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.DynamicRoutes = new(true)
			}),
		},
		{
			name: "disable system routes",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
				f[systemRoutesFlag] = "false"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.SystemRoutes = new(false)
			}),
		},
		{
			name: "enable system routes",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
				f[systemRoutesFlag] = "true"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.SystemRoutes = new(true)
			}),
		},
		{
			name: "combined updates",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[descriptionFlag] = "my routing table"
				f[dynamicRoutesFlag] = "false"
				f[systemRoutesFlag] = "false"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Description = new("my routing table")
				m.DynamicRoutes = new(false)
				m.SystemRoutes = new(false)
			}),
		},
		{
			name: "no update flags",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
			}),
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
			name: "invalid labels",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[labelFlag] = "invalid"
			}),
		},
		{
			name: "invalid boolean",
			args: []string{testRoutingTableId},
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[dynamicRoutesFlag] = "invalid"
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
			testutils.TestParseInput(t, NewCmd, parseInput, tt.want, tt.args, tt.flagValues, tt.want != nil)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		name        string
		model       *inputModel
		wantPayload iaas.UpdateVPCRoutingTablePayload
	}{
		{
			name: "all fields",
			model: fixtureInputModel(func(m *inputModel) {
				m.Description = new("description")
				m.Labels = map[string]any{"env": "dev"}
				m.DynamicRoutes = new(true)
				m.SystemRoutes = new(true)
			}),
			wantPayload: fixturePayload(func(p *iaas.UpdateVPCRoutingTablePayload) {
				p.Description = new("description")
				p.Labels = map[string]any{"env": "dev"}
				p.DynamicRoutes = new(true)
				p.SystemRoutes = new(true)
			}),
		},
		{
			name: "empty labels",
			model: fixtureInputModel(func(m *inputModel) {
				m.Labels = map[string]any{}
			}),
			wantPayload: fixturePayload(func(p *iaas.UpdateVPCRoutingTablePayload) {
				p.Labels = map[string]any{}
			}),
		},
		{
			name: "omitted labels",
			model: fixtureInputModel(func(m *inputModel) {
				m.Description = new("description")
				m.DynamicRoutes = new(true)
				m.SystemRoutes = new(true)
			}),
			wantPayload: fixturePayload(func(p *iaas.UpdateVPCRoutingTablePayload) {
				p.Description = new("description")
				p.DynamicRoutes = new(true)
				p.SystemRoutes = new(true)
			}),
		},
		{
			name:        "name only leaves other fields unchanged",
			model:       fixtureInputModel(),
			wantPayload: fixturePayload(),
		},
		{
			name: "disabled routes",
			model: fixtureInputModel(func(m *inputModel) {
				m.DynamicRoutes = new(false)
				m.SystemRoutes = new(false)
			}),
			wantPayload: fixturePayload(func(p *iaas.UpdateVPCRoutingTablePayload) {
				p.DynamicRoutes = new(false)
				p.SystemRoutes = new(false)
			}),
		},
		{
			name: "clear description",
			model: fixtureInputModel(func(m *inputModel) {
				m.Description = new("")
			}),
			wantPayload: fixturePayload(func(p *iaas.UpdateVPCRoutingTablePayload) {
				p.Description = new("")
			}),
		},
		{
			name: "omitted name",
			model: fixtureInputModel(func(m *inputModel) {
				m.Name = nil
				m.Description = new("description")
			}),
			wantPayload: fixturePayload(func(p *iaas.UpdateVPCRoutingTablePayload) {
				p.Name = nil
				p.Description = new("description")
			}),
		},
	}
	for _, tt := range tests {
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
