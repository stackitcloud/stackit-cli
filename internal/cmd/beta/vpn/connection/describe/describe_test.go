package describe

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "test")
	testProjectId    = uuid.NewString()
	testGatewayID    = uuid.NewString()
	testConnectionID = uuid.NewString()
	testClient, _    = vpn.NewAPIClient(
		sdkConfig.WithoutAuthentication(),
	)
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testConnectionID,
	}
	for _, m := range mods {
		m(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		gatewayIdFlag:             testGatewayID,
	}
	for _, m := range mods {
		m(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
		},
		GatewayId:    utils.Ptr(testGatewayID),
		ConnectionId: testConnectionID,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiGetGatewayConnectionRequest)) vpn.ApiGetGatewayConnectionRequest {
	request := testClient.DefaultAPI.GetGatewayConnection(testCtx, testProjectId, "", testGatewayID, testConnectionID)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureResponse(mods ...func(resp *vpn.ConnectionResponse)) *vpn.ConnectionResponse {
	resp := &vpn.ConnectionResponse{
		Id:          utils.Ptr(testConnectionID),
		DisplayName: "test-connection",
		Enabled:     utils.Ptr(true),
		Labels: &map[string]string{
			"env": "prod",
		},
		LocalSubnets:  []string{"10.0.0.0/24"},
		RemoteSubnets: []string{"192.168.0.0/24"},
		StaticRoutes:  []string{"10.1.0.0/24"},
		Tunnel1: vpn.TunnelConfiguration{
			RemoteAddress: "1.2.3.4",
			Bgp: &vpn.BGPTunnelConfig{
				RemoteAsn: 65000,
			},
			Peering: &vpn.PeeringConfig{
				LocalAddress:  utils.Ptr("169.254.0.1"),
				RemoteAddress: utils.Ptr("169.254.0.2"),
			},
			Phase1: vpn.TunnelConfigurationPhase1{
				DhGroups:             []vpn.PhaseDhGroupsInner{"14"},
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
				RekeyTime:            utils.Ptr(int32(3600)),
			},
			Phase2: vpn.TunnelConfigurationPhase2{
				DhGroups:             []vpn.PhaseDhGroupsInner{"14"},
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
				RekeyTime:            utils.Ptr(int32(3600)),
				DpdAction:            utils.Ptr(vpn.TunnelConfigurationPhase2AllOfDpdAction("restart")),
				StartAction:          utils.Ptr(vpn.TunnelConfigurationPhase2AllOfStartAction("start")),
			},
		},
		Tunnel2: vpn.TunnelConfiguration{
			RemoteAddress: "5.6.7.8",
			Phase1: vpn.TunnelConfigurationPhase1{
				DhGroups:             []vpn.PhaseDhGroupsInner{"14"},
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
				RekeyTime:            utils.Ptr(int32(3600)),
			},
			Phase2: vpn.TunnelConfigurationPhase2{
				DhGroups:             []vpn.PhaseDhGroupsInner{"14"},
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
				RekeyTime:            utils.Ptr(int32(3600)),
				DpdAction:            utils.Ptr(vpn.TunnelConfigurationPhase2AllOfDpdAction("restart")),
				StartAction:          utils.Ptr(vpn.TunnelConfigurationPhase2AllOfStartAction("start")),
			},
		},
	}

	for _, mod := range mods {
		mod(resp)
	}
	return resp
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
			description: "no args",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no gateway id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, gatewayIdFlag)
			}),
			isValid: false,
		},
		{
			description: "no project id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
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
		description    string
		model          *inputModel
		expectedResult vpn.ApiGetGatewayConnectionRequest
	}{
		{
			description:    "base",
			model:          fixtureInputModel(),
			expectedResult: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedResult,
				cmp.AllowUnexported(tt.expectedResult),
				cmpopts.IgnoreUnexported(vpn.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		resp        *vpn.ConnectionResponse
		wantErr     bool
	}{
		{
			description: "nil response",
			model:       fixtureInputModel(),
			resp:        nil,
			wantErr:     true,
		},
		{
			description: "full response",
			model:       fixtureInputModel(),
			resp:        fixtureResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			err := outputResult(params.Printer, tt.model, tt.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
