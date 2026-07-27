package status

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &vpn.APIClient{DefaultAPI: &vpn.DefaultAPIService{}}

var testProjectId = uuid.NewString()
var testRegion = "eu01"

var testGatewayId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testGatewayId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		GatewayId: testGatewayId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiGetGatewayStatusRequest)) vpn.ApiGetGatewayStatusRequest {
	request := testClient.DefaultAPI.GetGatewayStatus(testCtx, testProjectId, testRegion, testGatewayId)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "gateway id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "gateway id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest vpn.ApiGetGatewayStatusRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, vpn.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		gatewayId    string
		projectLabel string
		gateway      *vpn.GatewayStatusResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "set empty response",
			args: args{
				gateway: &vpn.GatewayStatusResponse{},
			},
			wantErr: false,
		},
		{
			name: "set response",
			args: args{
				gateway: &vpn.GatewayStatusResponse{
					Id:            utils.Ptr(testGatewayId),
					Connections:   []vpn.ConnectionStatusResponse{},
					DisplayName:   utils.Ptr("test"),
					ErrorMessage:  nil,
					GatewayStatus: utils.Ptr(vpn.GATEWAYSTATUS_READY),
					Tunnels: []vpn.VPNTunnels{
						{
							BgpStatus: *vpn.NewNullableBGPStatus(&vpn.BGPStatus{
								Peers: []vpn.BGPStatusPeers{
									{
										LocalAs:    23,
										PeerUptime: "~10s",
										PfxRcd:     3,
										PfxSnt:     5,
										RemoteAs:   224,
										RemoteIP:   "1.1.1.1",
										State:      "Healthy",
									},
									{
										LocalAs:    25,
										PeerUptime: "~14s",
										PfxRcd:     1,
										PfxSnt:     99,
										RemoteAs:   4,
										RemoteIP:   "2.2.2.2",
										State:      "Unhealthy",
									},
								},
								Routes: []vpn.BGPStatusRoutes{
									{
										Network: "home",
										Origin:  "root",
										Path:    "~/",
										PeerId:  "5",
										Weight:  1,
									},
									{
										Network: "remote",
										Origin:  "internet",
										Path:    "/path/",
										PeerId:  "3",
										Weight:  44,
									},
								},
							}),
							InstanceState:     utils.Ptr(vpn.GATEWAYSTATUS_READY),
							InternalNextHopIP: utils.Ptr("1.2.3.4"),
							Name:              utils.Ptr(vpn.VPNTUNNELSNAME_TUNNEL1),
							PublicIP:          utils.Ptr("9.9.9.9"),
						},
						{
							BgpStatus: *vpn.NewNullableBGPStatus(&vpn.BGPStatus{
								Peers: []vpn.BGPStatusPeers{
									{
										LocalAs:    23,
										PeerUptime: "~10s",
										PfxRcd:     3,
										PfxSnt:     5,
										RemoteAs:   224,
										RemoteIP:   "1.1.1.1",
										State:      "Healthy",
									},
									{
										LocalAs:    25,
										PeerUptime: "~14s",
										PfxRcd:     1,
										PfxSnt:     99,
										RemoteAs:   4,
										RemoteIP:   "2.2.2.2",
										State:      "Unhealthy",
									},
								},
								Routes: []vpn.BGPStatusRoutes{
									{
										Network: "home",
										Origin:  "root",
										Path:    "~/",
										PeerId:  "5",
										Weight:  1,
									},
									{
										Network: "remote",
										Origin:  "internet",
										Path:    "/path/",
										PeerId:  "3",
										Weight:  44,
									},
								},
							}),
							InstanceState:     utils.Ptr(vpn.GATEWAYSTATUS_PENDING),
							InternalNextHopIP: utils.Ptr("4.4.4.4"),
							Name:              utils.Ptr(vpn.VPNTUNNELSNAME_TUNNEL2),
							PublicIP:          utils.Ptr("3.3.3.3"),
						},
					},
				},
			},
		},
	}
	printer := print.NewPrinter(
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(printer, tt.args.outputFormat, tt.args.projectLabel, tt.args.gatewayId, tt.args.gateway); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
