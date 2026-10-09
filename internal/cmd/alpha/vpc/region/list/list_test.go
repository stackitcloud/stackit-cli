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
			name:  "no region required",
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
		want        iaas.ApiListVPCRegionsRequest
	}{
		{
			description: "default",
			inputModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				VpcId: testVpcId,
			},
			want: client.ListVPCRegions(ctx, testProjectId, testVpcId),
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
	resp := &iaas.RegionalVPCList{
		Regions: map[string]iaas.RegionalVPC{
			"eu02": {
				Status: new("CREATED"),
				Ipv4: &iaas.RegionalVPCIPv4{
					DefaultNameservers: []string{"8.8.8.8"},
				},
			},
			"eu01": {},
		},
	}

	for _, tt := range []struct {
		name    string
		format  string
		resp    *iaas.RegionalVPCList
		want    []string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "empty",
			resp: &iaas.RegionalVPCList{},
			want: []string{`No regions found for VPC "my-vpc"`},
		},
		{
			name: "table",
			resp: resp,
			want: []string{"eu01", "eu02", "CREATED", "8.8.8.8"},
		},
		{
			name:   "json",
			format: print.JSONOutputFormat,
			resp:   resp,
			want:   []string{`"regions"`, `"eu02"`, `"status": "CREATED"`},
		},
		{
			name:   "yaml",
			format: print.YAMLOutputFormat,
			resp:   resp,
			want:   []string{"regions:", "eu02:", "status: CREATED"},
		},
		{
			name:   "empty json",
			format: print.JSONOutputFormat,
			resp: &iaas.RegionalVPCList{
				Regions: map[string]iaas.RegionalVPC{},
			},
			want: []string{`"regions": {}`},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()

			err := outputResult(params.Printer, tt.format, "my-vpc", tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, want := range tt.want {
				if !strings.Contains(params.Out.String(), want) {
					t.Errorf("output %q does not contain %q", params.Out.String(), want)
				}
			}

			if tt.name == "table" && strings.Index(params.Out.String(), "eu01") > strings.Index(params.Out.String(), "eu02") {
				t.Errorf("regions are not sorted: %s", params.Out.String())
			}
		})
	}
}
