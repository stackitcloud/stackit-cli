package describe

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
					VpcId: testVpcId,
				}
			}

			testutils.TestParseInput(t, NewCmd, parseInput, expected, tt.args, flagValues, tt.valid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	ctx := context.Background()
	client := &iaas.DefaultAPIServiceMock{}

	tests := []struct {
		description string
		inputModel  *inputModel
		want        iaas.ApiGetVPCRegionRequest
	}{
		{
			description: "default",
			inputModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    "eu02",
				},
				VpcId: testVpcId,
			},
			want: client.GetVPCRegion(ctx, testProjectId, testVpcId, "eu02"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := buildRequest(ctx, tt.inputModel, client)
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(tt.want), cmpopts.EquateComparable(ctx, iaas.DefaultAPIService{})); diff != "" {
				t.Fatalf("request mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	resp := &iaas.RegionalVPC{
		Status: new("CREATED"),
		Ipv4: &iaas.RegionalVPCIPv4{
			DefaultNameservers: []string{"8.8.8.8", "8.8.4.4"},
		},
	}

	for _, tt := range []struct {
		name    string
		format  string
		vpcName string
		resp    *iaas.RegionalVPC
		want    []string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "empty",
			resp: &iaas.RegionalVPC{},
			want: []string{testVpcId, "eu02"},
		},
		{
			name:    "table",
			vpcName: "my-vpc",
			resp:    resp,
			want:    []string{testVpcId, "my-vpc", "eu02", "CREATED", "8.8.8.8,8.8.4.4"},
		},
		{
			name:   "json",
			format: print.JSONOutputFormat,
			resp:   resp,
			want:   []string{`"status": "CREATED"`, `"defaultNameservers"`},
		},
		{
			name:   "yaml",
			format: print.YAMLOutputFormat,
			resp:   resp,
			want:   []string{"status: CREATED", "defaultNameservers:"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.format, "eu02", testVpcId, tt.vpcName, tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, want := range tt.want {
				if !strings.Contains(params.Out.String(), want) {
					t.Errorf("output %q does not contain %q", params.Out.String(), want)
				}
			}

			if tt.format == "" && tt.vpcName == "" && strings.Contains(params.Out.String(), "NAME") {
				t.Errorf("unexpected name row: %s", params.Out.String())
			}
		})
	}
}
