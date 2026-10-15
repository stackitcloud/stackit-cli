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
	testProjectId = uuid.NewString()
	testVpcId     = uuid.NewString()
)

func TestParseInput(t *testing.T) {
	for _, tt := range []struct {
		name        string
		flag        string
		value       string
		remove      bool
		args        []string
		valid       bool
		nameservers []string
	}{
		{
			name:  "defaults",
			valid: true,
		},
		{
			name:        "nameservers",
			flag:        ipv4DefaultNameserversFlag,
			value:       "8.8.8.8,8.8.4.4",
			valid:       true,
			nameservers: []string{"8.8.8.8", "8.8.4.4"},
		},
		{
			name:        "empty nameservers",
			flag:        ipv4DefaultNameserversFlag,
			valid:       true,
			nameservers: []string{},
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

			var expected *inputModel
			if tt.valid {
				expected = &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{
						ProjectId: testProjectId,
						Region:    "eu02",
						Verbosity: globalflags.VerbosityDefault,
					},
					VpcId:                  testVpcId,
					IPv4DefaultNameservers: tt.nameservers,
				}
			}

			testutils.TestParseInput(t, NewCmd, parseInput, expected, tt.args, flagValues, tt.valid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	buildRequestInputModelDefault := func(opts ...func(m *inputModel)) *inputModel {
		model := &inputModel{
			GlobalFlagModel: &globalflags.GlobalFlagModel{
				ProjectId: testProjectId,
				Region:    "eu02",
			},
			VpcId:                  testVpcId,
			IPv4DefaultNameservers: []string{"8.8.8.8", "8.8.4.4"},
		}

		for _, opt := range opts {
			opt(model)
		}

		return model
	}

	buildRequestPayloadDefault := func(opts ...func(p *iaas.CreateVPCRegionPayload)) iaas.CreateVPCRegionPayload {
		payload := iaas.CreateVPCRegionPayload{
			Ipv4: &iaas.RegionalVPCIPv4{
				DefaultNameservers: []string{"8.8.8.8", "8.8.4.4"},
			},
		}

		for _, opt := range opts {
			opt(&payload)
		}

		return payload
	}

	tests := []struct {
		description string
		inputModel  *inputModel
		wantPayload iaas.CreateVPCRegionPayload
	}{
		{
			description: "default",
			inputModel:  buildRequestInputModelDefault(),
			wantPayload: buildRequestPayloadDefault(),
		},
		{
			description: "nameservers is empty",
			inputModel: buildRequestInputModelDefault(func(m *inputModel) {
				m.IPv4DefaultNameservers = []string{}
			}),
			wantPayload: buildRequestPayloadDefault(func(p *iaas.CreateVPCRegionPayload) {
				p.Ipv4.DefaultNameservers = []string{}
			}),
		},
		{
			description: "nameservers is nil",
			inputModel: buildRequestInputModelDefault(func(m *inputModel) {
				m.IPv4DefaultNameservers = nil
			}),
			wantPayload: buildRequestPayloadDefault(func(p *iaas.CreateVPCRegionPayload) {
				p.Ipv4 = nil
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := context.Background()
			client := &iaas.DefaultAPIServiceMock{}

			want := client.CreateVPCRegion(ctx, testProjectId, testVpcId, "eu02").CreateVPCRegionPayload(tt.wantPayload)
			got := buildRequest(ctx, tt.inputModel, client)
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
		async   bool
		resp    *iaas.RegionalVPC
		want    string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "created",
			resp: &iaas.RegionalVPC{},
			want: "Created region configuration",
		},
		{
			name:  "async",
			async: true,
			resp:  &iaas.RegionalVPC{},
			want:  "Triggered creation of",
		},
		{
			name:   "json",
			format: print.JSONOutputFormat,
			resp: &iaas.RegionalVPC{
				Status: new("CREATED"),
			},
			want: `"status": "CREATED"`,
		},
		{
			name:   "yaml",
			format: print.YAMLOutputFormat,
			resp: &iaas.RegionalVPC{
				Status: new("CREATED"),
			},
			want: "status: CREATED",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.format, tt.async, "eu02", "my-vpc", tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !strings.Contains(params.Out.String(), tt.want) {
				t.Errorf("output %q does not contain %q", params.Out.String(), tt.want)
			}
		})
	}
}
