package list

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const testRegion = "eu01"

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRoutingTableId = uuid.NewString()

const testLabelSelectorFlag = "key1=value1,key2=value2"

var testLabels = &map[string]string{
	"key1": "value1",
	"key2": "value2",
}

var testLimitFlag = int64(10)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,
		organizationIdFlag:     testOrgId,
		networkAreaIdFlag:      testNetworkAreaId,
		routingTableIdFlag:     testRoutingTableId,
		labelSelectorFlag:      testLabelSelectorFlag,
		limitFlag:              strconv.Itoa(int(testLimitFlag)),
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
		LabelSelector:  utils.Ptr(testLabelSelectorFlag),
		Limit:          utils.Ptr(testLimitFlag),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "routing-table-id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, routingTableIdFlag)
			}),
			isValid: false,
		},
		{
			description: "network-area-id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "org-id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "labels missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelSelectorFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.LabelSelector = nil
			}),
		},
		{
			description: "limit missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, limitFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = nil
			}),
		},
		{
			description: "invalid limit flag",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "negative limit flag",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "-10"
			}),
			isValid: false,
		},
		{
			description: "limit zero flag",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
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
		routes       []iaas.Route
		wantErr      bool
	}{
		{
			name:         "nil routes should return error",
			outputFormat: print.PrettyOutputFormat,
			routes:       nil,
			wantErr:      true,
		},
		{
			name:         "empty routes list",
			outputFormat: print.PrettyOutputFormat,
			routes:       []iaas.Route{},
			wantErr:      false,
		},
		{
			name:         "empty routes list json output",
			outputFormat: print.JSONOutputFormat,
			routes:       []iaas.Route{},
			wantErr:      false,
		},
		{
			name:         "empty routes list json output",
			outputFormat: print.YAMLOutputFormat,
			routes:       []iaas.Route{},
			wantErr:      false,
		},
		{
			name:         "route list with empty struct",
			outputFormat: print.PrettyOutputFormat,
			routes:       []iaas.Route{{}},
			wantErr:      false,
		},
		{
			name:         "pretty output with one route",
			outputFormat: print.PrettyOutputFormat,
			routes:       []iaas.Route{dummyRoute},
			wantErr:      false,
		},
		{
			name:         "pretty output with multiple routes",
			outputFormat: print.PrettyOutputFormat,
			routes:       []iaas.Route{dummyRoute, dummyRoute, dummyRoute},
			wantErr:      false,
		},
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			routes:       []iaas.Route{dummyRoute},
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			routes:       []iaas.Route{dummyRoute},
			wantErr:      false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.outputFormat, tt.routes); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
