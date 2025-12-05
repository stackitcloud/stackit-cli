package update

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

const testRegion = "eu01"

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRoutingTableId = uuid.NewString()
var testRouteId = uuid.NewString()

const testLabelSelectorFlag = "key1=value1,key2=value2"

var testLabels = &map[string]string{
	"key1": "value1",
	"key2": "value2",
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,
		organizationIdFlag:     testOrgId,
		networkAreaIdFlag:      testNetworkAreaId,
		routingTableIdFlag:     testRoutingTableId,
		labelFlag:              testLabelSelectorFlag,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		OrganizationId: testOrgId,
		NetworkAreaId:  testNetworkAreaId,
		RoutingTableId: testRoutingTableId,
		RouteId:        testRouteId,
		Labels:         testLabels,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testRouteId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureRequest(mods ...func(req *iaas.ApiUpdateRouteOfRoutingTableRequest)) iaas.ApiUpdateRouteOfRoutingTableRequest {
	req := testClient.UpdateRouteOfRoutingTable(
		testCtx,
		testOrgId,
		testNetworkAreaId,
		testRegion,
		testRoutingTableId,
		testRouteId,
	)

	payload := iaas.UpdateRouteOfRoutingTablePayload{
		Labels: utils.ConvertStringMapToInterfaceMap(testLabels),
	}

	req = req.UpdateRouteOfRoutingTablePayload(payload)

	for _, mod := range mods {
		mod(&req)
	}

	return req
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		argValues     []string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			argValues:     fixtureArgValues(),
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
			description: "routing-table-id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, routingTableIdFlag)
			}),
			isValid: false,
		},
		{
			description: "network-area-id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "org-id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "routing-table-id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, routingTableIdFlag)
			}),
			isValid: false,
		},
		{
			description:   "arg value missing",
			argValues:     []string{""},
			flagValues:    fixtureFlagValues(),
			isValid:       false,
			expectedModel: fixtureInputModel(),
		},
		{
			description:   "arg value wrong",
			argValues:     []string{"foo-bar"},
			flagValues:    fixtureFlagValues(),
			isValid:       false,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "labels are missing",
			argValues:   []string{},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid label format",
			argValues:   []string{},
			flagValues:  map[string]string{labelFlag: "invalid-label"},
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
		expectedRequest iaas.ApiUpdateRouteOfRoutingTableRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "labels nil",
			model: fixtureInputModel(func(m *inputModel) {
				m.Labels = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateRouteOfRoutingTableRequest) {
				*request = (*request).UpdateRouteOfRoutingTablePayload(iaas.UpdateRouteOfRoutingTablePayload{
					Labels: nil,
				})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			gotReq := buildRequest(testCtx, tt.model, testClient)

			if diff := cmp.Diff(
				tt.expectedRequest,
				gotReq,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			); diff != "" {
				t.Errorf("buildRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	dummyRoute := iaas.Route{
		Id: utils.Ptr("route-foo"),
		Destination: &iaas.RouteDestination{
			DestinationCIDRv4: &iaas.DestinationCIDRv4{
				Type:  utils.Ptr("cidrv4"),
				Value: utils.Ptr("10.0.0.0/24"),
			},
		},
		Nexthop: &iaas.RouteNexthop{
			NexthopIPv4: &iaas.NexthopIPv4{
				Type:  utils.Ptr("ipv4"),
				Value: utils.Ptr("10.0.0.1"),
			},
		},
		Labels:    utils.ConvertStringMapToInterfaceMap(testLabels),
		CreatedAt: utils.Ptr(time.Now()),
		UpdatedAt: utils.Ptr(time.Now()),
	}

	tests := []struct {
		name         string
		outputFormat string
		route        *iaas.Route
		wantErr      bool
	}{
		{
			name:         "nil route should return error",
			outputFormat: print.PrettyOutputFormat,
			route:        nil,
			wantErr:      true,
		},
		{
			name:         "empty route",
			outputFormat: print.PrettyOutputFormat,
			route:        &iaas.Route{},
			// should fail on pretty format
			wantErr: true,
		},
		{
			name:         "pretty output with one route",
			outputFormat: print.PrettyOutputFormat,
			route:        &dummyRoute,
			wantErr:      false,
		},
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			route:        &dummyRoute,
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			route:        &dummyRoute,
			wantErr:      false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.outputFormat, "", "", tt.route); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
