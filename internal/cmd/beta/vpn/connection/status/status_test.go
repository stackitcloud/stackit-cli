package status

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

func fixtureRequest(mods ...func(request *vpn.ApiGetGatewayConnectionStatusRequest)) vpn.ApiGetGatewayConnectionStatusRequest {
	request := testClient.DefaultAPI.GetGatewayConnectionStatus(testCtx, testProjectId, "", testGatewayID, testConnectionID)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureResponse(mods ...func(resp *vpn.ConnectionStatusResponse)) *vpn.ConnectionStatusResponse {
	resp := &vpn.ConnectionStatusResponse{
		Id:          utils.Ptr(testConnectionID),
		DisplayName: utils.Ptr("test-connection"),
		Enabled:     utils.Ptr(true),
		Tunnels: []vpn.TunnelStatus{
			{
				Name:        utils.Ptr(vpn.TunnelStatusName("tunnel1")),
				Established: utils.Ptr(true),
				Phase1: &vpn.Phase1Status{
					DhGroup:             utils.Ptr("MODP2048"),
					EncryptionAlgorithm: utils.Ptr("AES_GCM_16"),
					IntegrityAlgorithm:  utils.Ptr("SHA_256"),
					State:               utils.Ptr("INSTALLED"),
				},
				Phase2: &vpn.Phase2Status{
					BytesIn:             utils.Ptr("453533"),
					BytesOut:            utils.Ptr("46459064"),
					DhGroup:             utils.Ptr("MODP2048"),
					Encap:               utils.Ptr("yes"),
					EncryptionAlgorithm: utils.Ptr("AES_GCM_16"),
					IntegrityAlgorithm:  utils.Ptr("SHA_256"),
					PacketsIn:           utils.Ptr("1534134"),
					PacketsOut:          utils.Ptr("65847343"),
					Protocol:            utils.Ptr("ESP"),
					State:               utils.Ptr("ESTABLISHED"),
				},
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
		name     string
		model    *inputModel
		expected vpn.ApiGetGatewayConnectionStatusRequest
	}{
		{
			name:     "base",
			model:    fixtureInputModel(),
			expected: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			diff := cmp.Diff(request, tt.expected,
				cmp.AllowUnexported(tt.expected),
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
		resp        *vpn.ConnectionStatusResponse
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
