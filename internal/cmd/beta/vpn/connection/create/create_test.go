package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/spf13/cobra"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "test")
	testProjectId = uuid.NewString()
	testGatewayID = uuid.NewString()
	testClient, _ = vpn.NewAPIClient(
		sdkConfig.WithoutAuthentication(),
	)
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag:                    testProjectId,
		gatewayIdFlag:                                testGatewayID,
		displayNameFlag:                              "test-connection",
		tunnel1RemoteAddressFlag:                     "1.2.3.4",
		tunnel1PreSharedKeyFlag:                      "test-psk-1",
		tunnel1Phase1EncryptionAlgorithmsFlag.Name(): "aes256",
		tunnel1Phase1IntegrityAlgorithmsFlag.Name():  "sha2_256",
		tunnel1Phase2EncryptionAlgorithmsFlag.Name(): "aes256",
		tunnel1Phase2IntegrityAlgorithmsFlag.Name():  "sha2_256",
		tunnel2RemoteAddressFlag:                     "5.6.7.8",
		tunnel2PreSharedKeyFlag:                      "test-psk-2",
		tunnel2Phase1EncryptionAlgorithmsFlag.Name(): "aes256",
		tunnel2Phase1IntegrityAlgorithmsFlag.Name():  "sha2_256",
		tunnel2Phase2EncryptionAlgorithmsFlag.Name(): "aes256",
		tunnel2Phase2IntegrityAlgorithmsFlag.Name():  "sha2_256",
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
		GatewayId:   testGatewayID,
		DisplayName: "test-connection",
		Enabled:     nil,
		Tunnel1: tunnelInputModel{
			RemoteAddress:              "1.2.3.4",
			PreSharedKey:               "test-psk-1",
			Phase1EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
			Phase1IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
			Phase2EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
			Phase2IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
		},
		Tunnel2: tunnelInputModel{
			RemoteAddress:              "5.6.7.8",
			PreSharedKey:               "test-psk-2",
			Phase1EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
			Phase1IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
			Phase2EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
			Phase2IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
		},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiCreateGatewayConnectionRequest)) vpn.ApiCreateGatewayConnectionRequest {
	request := testClient.DefaultAPI.CreateGatewayConnection(testCtx, testProjectId, "", testGatewayID)
	payload := vpn.CreateGatewayConnectionPayload{
		DisplayName: "test-connection",
		Enabled:     nil,
		Tunnel1: vpn.TunnelConfiguration{
			RemoteAddress: "1.2.3.4",
			PreSharedKey:  utils.Ptr("test-psk-1"),
			Phase1: vpn.TunnelConfigurationPhase1{
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
			},
			Phase2: vpn.TunnelConfigurationPhase2{
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
			},
		},
		Tunnel2: vpn.TunnelConfiguration{
			RemoteAddress: "5.6.7.8",
			PreSharedKey:  utils.Ptr("test-psk-2"),
			Phase1: vpn.TunnelConfigurationPhase1{
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
			},
			Phase2: vpn.TunnelConfigurationPhase2{
				EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
				IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
			},
		},
	}
	request = request.CreateGatewayConnectionPayload(payload)
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
			argValues:     []string{},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no flags",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "missing project id",
			argValues:   []string{},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "missing gateway id",
			argValues:   []string{},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, gatewayIdFlag)
			}),
			isValid: false,
		},
		{
			description: "missing display name",
			argValues:   []string{},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "missing tunnel1 remote address",
			argValues:   []string{},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, tunnel1RemoteAddressFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, func(printer *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
				return parseInput(printer, cmd)
			}, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description    string
		model          *inputModel
		expectedResult vpn.ApiCreateGatewayConnectionRequest
	}{
		{
			description:    "base",
			model:          fixtureInputModel(),
			expectedResult: fixtureRequest(),
		},
		{
			description: "with optional fields",
			model: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{"env": "prod"}
				model.LocalSubnets = []string{"10.0.0.0/24"}
				model.RemoteSubnets = []string{"192.168.0.0/24"}
				model.StaticRoutes = []string{"10.1.0.0/24"}
				model.Tunnel1.BgpRemoteAsn = utils.Ptr(int64(65000))
				model.Tunnel1.PeeringLocalAddress = utils.Ptr("169.254.0.1")
				model.Tunnel1.PeeringRemoteAddress = utils.Ptr("169.254.0.2")
				model.Tunnel1.Phase1DhGroups = []vpn.PhaseDhGroupsInner{"14"}
				model.Tunnel1.Phase1RekeyTime = utils.Ptr(int32(3600))
				model.Tunnel1.Phase2DhGroups = []vpn.PhaseDhGroupsInner{"14"}
				model.Tunnel1.Phase2RekeyTime = utils.Ptr(int32(3600))
				model.Tunnel1.Phase2DpdAction = utils.Ptr(vpn.TunnelConfigurationPhase2AllOfDpdAction("restart"))
				model.Tunnel1.Phase2StartAction = utils.Ptr(vpn.TunnelConfigurationPhase2AllOfStartAction("start"))
			}),
			expectedResult: fixtureRequest(func(request *vpn.ApiCreateGatewayConnectionRequest) {
				payload := vpn.CreateGatewayConnectionPayload{
					DisplayName:   "test-connection",
					Enabled:       nil,
					Labels:        &map[string]string{"env": "prod"},
					LocalSubnets:  []string{"10.0.0.0/24"},
					RemoteSubnets: []string{"192.168.0.0/24"},
					StaticRoutes:  []string{"10.1.0.0/24"},
					Tunnel1: vpn.TunnelConfiguration{
						RemoteAddress: "1.2.3.4",
						PreSharedKey:  utils.Ptr("test-psk-1"),
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
						PreSharedKey:  utils.Ptr("test-psk-2"),
						Phase1: vpn.TunnelConfigurationPhase1{
							EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
							IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
						},
						Phase2: vpn.TunnelConfigurationPhase2{
							EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
							IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
						},
					},
				}
				*request = request.CreateGatewayConnectionPayload(payload)
			}),
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
		expected    string
		wantErr     bool
	}{
		{
			description: "nil response",
			model:       fixtureInputModel(),
			resp:        nil,
			wantErr:     true,
			expected:    "",
		},
		{
			description: "success",
			model:       fixtureInputModel(),
			resp: &vpn.ConnectionResponse{
				Id: utils.Ptr("conn-1234"),
			},
			expected: fmt.Sprintf("Created VPN connection \"conn-1234\" for gateway %q in project %q.\n", testGatewayID, testProjectId),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			err := outputResult(params.Printer, tt.model, testProjectId, tt.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && params.Out.String() != tt.expected {
				t.Errorf("want:\n%s\ngot:\n%s", tt.expected, params.Out.String())
			}
		})
	}
}
