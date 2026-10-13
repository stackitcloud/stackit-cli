package create

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
		VpcId:         testVpcId,
		Name:          "rt",
		DynamicRoutes: true,
		SystemRoutes:  true,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixturePayload(mods ...func(*iaas.AddVPCRoutingTablePayload)) iaas.AddVPCRoutingTablePayload {
	payload := iaas.AddVPCRoutingTablePayload{
		Name:          "rt",
		DynamicRoutes: new(true),
		SystemRoutes:  new(true),
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
			name:       "defaults",
			flagValues: fixtureFlagValues(),
			want:       fixtureInputModel(),
		},
		{
			name: "description",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[descriptionFlag] = "my routing table"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Description = new("my routing table")
			}),
		},
		{
			name: "empty description",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[descriptionFlag] = ""
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Description = new("")
			}),
		},
		{
			name: "labels",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[labelFlag] = "env=dev,team=network"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.Labels = map[string]any{
					"env":  "dev",
					"team": "network",
				}
			}),
		},
		{
			name: "empty labels",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[labelFlag] = ""
			}),
		},
		{
			name: "disable dynamic routes",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[dynamicRoutesFlag] = "false"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.DynamicRoutes = false
			}),
		},
		{
			name: "disable system routes",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[systemRoutesFlag] = "false"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.SystemRoutes = false
			}),
		},
		{
			name: "disable both route types",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[dynamicRoutesFlag] = "false"
				f[systemRoutesFlag] = "false"
			}),
			want: fixtureInputModel(func(m *inputModel) {
				m.DynamicRoutes = false
				m.SystemRoutes = false
			}),
		},
		{
			name: "missing name",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				delete(f, nameFlag)
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
			name: "invalid labels",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[labelFlag] = "invalid"
			}),
		},
		{
			name: "invalid boolean",
			flagValues: fixtureFlagValues(func(f map[string]string) {
				f[dynamicRoutesFlag] = "invalid"
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
	tests := []struct {
		name        string
		model       *inputModel
		wantPayload iaas.AddVPCRoutingTablePayload
	}{
		{
			name: "all fields",
			model: fixtureInputModel(func(m *inputModel) {
				m.Description = new("description")
				m.Labels = map[string]any{"env": "dev"}
			}),
			wantPayload: fixturePayload(func(p *iaas.AddVPCRoutingTablePayload) {
				p.Description = new("description")
				p.Labels = map[string]any{"env": "dev"}
			}),
		},
		{
			name: "empty labels",
			model: fixtureInputModel(func(m *inputModel) {
				m.Labels = map[string]any{}
			}),
			wantPayload: fixturePayload(func(p *iaas.AddVPCRoutingTablePayload) {
				p.Labels = map[string]any{}
			}),
		},
		{
			name:        "omitted optional fields",
			model:       fixtureInputModel(),
			wantPayload: fixturePayload(),
		},
		{
			name: "disabled routes",
			model: fixtureInputModel(func(m *inputModel) {
				m.DynamicRoutes = false
				m.SystemRoutes = false
			}),
			wantPayload: fixturePayload(func(p *iaas.AddVPCRoutingTablePayload) {
				p.DynamicRoutes = new(false)
				p.SystemRoutes = new(false)
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client := &iaas.DefaultAPIServiceMock{}

			want := client.AddVPCRoutingTable(ctx, testProjectId, testVpcId, "eu02").AddVPCRoutingTablePayload(tt.wantPayload)
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
			name: "created",
			resp: &iaas.VPCRoutingTable{
				Id: &testRoutingTableId,
			},
			want: "Created routing-table with ID \"" + testRoutingTableId + "\".",
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
