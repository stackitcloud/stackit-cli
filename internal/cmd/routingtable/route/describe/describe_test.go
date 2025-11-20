package describe

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const testRegion = "eu01"

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRoutingTableId = uuid.NewString()
var testRouteId = uuid.NewString()

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
		RouteID:        testRouteId,
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
			description: "invalid routing-table-id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[routingTableIdFlag] = "invalid-id"
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
			description: "invalid organization-id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-org"
			}),
			isValid: false,
		},
		{
			description: "invalid network-area-id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = "invalid-area"
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
		route        iaas.Route
		wantErr      bool
	}{
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			route:        dummyRoute,
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			route:        dummyRoute,
			wantErr:      false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.outputFormat, &tt.route); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
