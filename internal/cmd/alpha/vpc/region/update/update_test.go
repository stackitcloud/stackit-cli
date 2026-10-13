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
	testProjectId = uuid.NewString()
	testVpcId     = uuid.NewString()
)

func TestParseInput(t *testing.T) {
	for _, tt := range []struct {
		name   string
		flag   string
		value  string
		remove bool
		args   []string
		valid  bool
		clear  bool
	}{
		{
			name:  "nameservers",
			valid: true,
		},
		{
			name:  "clear nameservers",
			flag:  ipv4DefaultNameserversFlag,
			valid: true,
			clear: true,
		},
		{
			name:   "nothing to update",
			flag:   ipv4DefaultNameserversFlag,
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
			name: "unexpected argument",
			args: []string{"unexpected"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			flagValues := map[string]string{
				globalflags.ProjectIdFlag:  testProjectId,
				globalflags.RegionFlag:     "eu02",
				vpcIdFlag:                  testVpcId,
				ipv4DefaultNameserversFlag: "8.8.8.8,8.8.4.4",
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
					IPv4DefaultNameservers: []string{"8.8.8.8", "8.8.4.4"},
				}
				if tt.clear {
					expected.IPv4DefaultNameservers = []string{}
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

	buildRequestPayloadDefault := func(opts ...func(p *iaas.UpdateVPCRegionPayload)) iaas.UpdateVPCRegionPayload {
		payload := iaas.UpdateVPCRegionPayload{
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
		wantPayload iaas.UpdateVPCRegionPayload
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
			wantPayload: buildRequestPayloadDefault(func(p *iaas.UpdateVPCRegionPayload) {
				p.Ipv4.DefaultNameservers = []string{}
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := context.Background()
			client := &iaas.DefaultAPIServiceMock{}

			want := client.UpdateVPCRegion(ctx, testProjectId, testVpcId, "eu02").UpdateVPCRegionPayload(tt.wantPayload)
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
		resp    *iaas.RegionalVPC
		want    string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "updated",
			resp: &iaas.RegionalVPC{},
			want: "Updated region configuration",
		},
		{
			name:   "json",
			format: print.JSONOutputFormat,
			resp: &iaas.RegionalVPC{
				Status: new("UPDATED"),
			},
			want: `"status": "UPDATED"`,
		},
		{
			name:   "yaml",
			format: print.YAMLOutputFormat,
			resp: &iaas.RegionalVPC{
				Status: new("UPDATED"),
			},
			want: "status: UPDATED",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.format, "eu02", "my-vpc", tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !strings.Contains(params.Out.String(), tt.want) {
				t.Errorf("output %q does not contain %q", params.Out.String(), tt.want)
			}
		})
	}
}
