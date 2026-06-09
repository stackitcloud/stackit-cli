package list

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
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
		GatewayId: utils.Ptr(testGatewayID),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiListGatewayConnectionsRequest)) vpn.ApiListGatewayConnectionsRequest {
	request := testClient.DefaultAPI.ListGatewayConnections(testCtx, testProjectId, "", testGatewayID)
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
			description: "no project id",
			argValues:   []string{},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "no gateway id",
			argValues:   []string{},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, gatewayIdFlag)
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
		expectedResult vpn.ApiListGatewayConnectionsRequest
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
		resp        *vpn.ConnectionList
		wantErr     bool
	}{
		{
			description: "empty list",
			model:       fixtureInputModel(),
			resp: &vpn.ConnectionList{
				Connections: []vpn.ConnectionResponse{},
			},
		},
		{
			description: "nil response",
			model:       fixtureInputModel(),
			resp:        nil,
			wantErr:     true,
		},
		{
			description: "nil connections",
			model:       fixtureInputModel(),
			resp: &vpn.ConnectionList{
				Connections: nil,
			},
			wantErr: true,
		},
		{
			description: "with entries",
			model:       fixtureInputModel(),
			resp: &vpn.ConnectionList{
				Connections: []vpn.ConnectionResponse{
					{
						Id:          utils.Ptr("conn-1"),
						DisplayName: "test-conn-1",
						Enabled:     utils.Ptr(true),
						Labels: &map[string]string{
							"env": "prod",
						},
					},
					{
						Id:          utils.Ptr("conn-2"),
						DisplayName: "test-conn-2",
						Enabled:     utils.Ptr(false),
						Labels:      nil,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			err := outputResult(params.Printer, tt.model, testProjectId, tt.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
