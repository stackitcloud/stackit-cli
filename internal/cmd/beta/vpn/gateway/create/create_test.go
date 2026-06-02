package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

var projectIdFlag = globalflags.ProjectIdFlag
var regionFlag = globalflags.RegionFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &vpn.APIClient{DefaultAPI: &vpn.DefaultAPIService{}}

var testProjectId = uuid.NewString()
var testRegion = "eu01"

var testName = "test-name"
var testPlanId = "planId"
var testRoutingType = vpn.ROUTINGTYPE_POLICY_BASED
var testAvailabilityZoneTunnel1 = "eu01-1"
var testAvailabilityZoneTunnel2 = "eu01-2"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		regionFlag:    testRegion,

		availabilityZoneTunnel1Flag: testAvailabilityZoneTunnel1,
		availabilityZoneTunnel2Flag: testAvailabilityZoneTunnel2,
		nameFlag:                    testName,
		labelsFlag:                  "foo=bar",
		planIdFlag:                  testPlanId,
		routingTypeFlag:             string(testRoutingType),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		Name: testName,
		AvailabilityZone: vpn.CreateGatewayPayloadAvailabilityZones{
			Tunnel1: testAvailabilityZoneTunnel1,
			Tunnel2: testAvailabilityZoneTunnel2,
		},
		Bgp: nil,
		Labels: &map[string]string{
			"foo": "bar",
		},
		PlanId:      testPlanId,
		RoutingType: testRoutingType,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiCreateGatewayRequest)) vpn.ApiCreateGatewayRequest {
	request := testClient.DefaultAPI.CreateGateway(testCtx, testProjectId, vpn.Region(testRegion))
	request = request.CreateGatewayPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(request *vpn.CreateGatewayPayload)) vpn.CreateGatewayPayload {
	payload := vpn.CreateGatewayPayload{
		AvailabilityZones: vpn.CreateGatewayPayloadAvailabilityZones{
			Tunnel1: testAvailabilityZoneTunnel1,
			Tunnel2: testAvailabilityZoneTunnel2,
		},
		Bgp:         nil,
		DisplayName: testName,
		Labels: &map[string]string{
			"foo": "bar",
		},
		PlanId:      testPlanId,
		RoutingType: testRoutingType,
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelsFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
		},
		{
			description: "missing required name",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
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
		description     string
		model           *inputModel
		expectedRequest vpn.ApiCreateGatewayRequest
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
				cmp.AllowUnexported(vpn.NullableString{}),
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
		async        bool
		projectLabel string
		item         *vpn.GatewayResponse
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
				item: &vpn.GatewayResponse{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.async, tt.args.projectLabel, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
