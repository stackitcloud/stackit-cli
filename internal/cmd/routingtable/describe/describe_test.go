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
	"github.com/stackitcloud/stackit-sdk-go/services/iaasalpha"
)

const testRegion = "eu01"

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRoutingTableId = uuid.NewString()

var testLabels = &map[string]string{
	"key1": "value1",
	"key2": "value2",
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,
		organizationIdFlag:     testOrgId,
		networkAreaIdFlag:      testNetworkAreaId,
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
		OrganizationId: utils.Ptr(testOrgId),
		NetworkAreaId:  utils.Ptr(testNetworkAreaId),
		RoutingTableId: utils.Ptr(testRoutingTableId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testRoutingTableId,
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestOutputResult(t *testing.T) {
	dummyRouteTable := iaasalpha.RoutingTable{
		CreatedAt:    utils.Ptr(time.Now()),
		Default:      nil,
		Description:  utils.Ptr("description"),
		Id:           utils.Ptr("route-foo"),
		Labels:       utils.ConvertStringMapToInterfaceMap(testLabels),
		Name:         utils.Ptr("route-foo"),
		SystemRoutes: utils.Ptr(true),
		UpdatedAt:    utils.Ptr(time.Now()),
	}

	tests := []struct {
		name         string
		outputFormat string
		routingTable iaasalpha.RoutingTable
		wantErr      bool
	}{
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			routingTable: dummyRouteTable,
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			routingTable: dummyRouteTable,
			wantErr:      false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.outputFormat, &tt.routingTable); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
