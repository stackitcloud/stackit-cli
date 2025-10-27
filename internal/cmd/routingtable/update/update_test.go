package update

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	pprint "github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const testRegion = "eu01"

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRoutingTableId = uuid.NewString()

const testRoutingTableName = "test"
const testRoutingTableDescription = "test"
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
		descriptionFlag:        testRoutingTableDescription,
		nameFlag:               testRoutingTableName,
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
		OrganizationId: utils.Ptr(testOrgId),
		NetworkAreaId:  utils.Ptr(testNetworkAreaId),
		Name:           utils.Ptr(testRoutingTableName),
		Description:    utils.Ptr(testRoutingTableDescription),
		Labels:         utils.Ptr(*testLabels),
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
			description: "base",
			flagValues:  fixtureFlagValues(),
			argValues:   fixtureArgValues(),
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = &testRoutingTableId
			}),
		},
		{
			description: "dynamic_routes disabled",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nonDynamicRoutesFlag] = "true"
			}),
			argValues: fixtureArgValues(),
			isValid:   true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NonDynamicRoutes = true
				model.RoutingTableId = &testRoutingTableId
			}),
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

func TestOutputResult(t *testing.T) {
	dummyRoutingTable := iaas.RoutingTable{
		Id:            utils.Ptr("id-foo"),
		Name:          utils.Ptr("route-table-foo"),
		Description:   utils.Ptr("description-foo"),
		SystemRoutes:  utils.Ptr(true),
		DynamicRoutes: utils.Ptr(true),
		Labels:        utils.ConvertStringMapToInterfaceMap(testLabels),
		CreatedAt:     utils.Ptr(time.Now()),
		UpdatedAt:     utils.Ptr(time.Now()),
	}

	tests := []struct {
		name         string
		outputFormat string
		routingTable *iaas.RoutingTable
		wantErr      bool
	}{
		{
			name:         "nil routing-table should return error",
			outputFormat: "",
			routingTable: nil,
			wantErr:      true,
		},
		{
			name:         "empty routing-table",
			outputFormat: "",
			routingTable: &iaas.RoutingTable{},
			wantErr:      true,
		},
		{
			name:         "table output routing-table",
			outputFormat: "",
			routingTable: &dummyRoutingTable,
			wantErr:      false,
		},
		{
			name:         "json output routing-table",
			outputFormat: pprint.JSONOutputFormat,
			routingTable: &dummyRoutingTable,
			wantErr:      false,
		},
		{
			name:         "yaml output routing-table",
			outputFormat: pprint.YAMLOutputFormat,
			routingTable: &dummyRoutingTable,
			wantErr:      false,
		},
	}

	p := pprint.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.outputFormat, "network-area-id", tt.routingTable); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
