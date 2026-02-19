package update

import (
	"context"
	"strconv"
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

const testRoutingTableName = "test"
const testRoutingTableDescription = "test"
const testLabelSelectorFlag = "key1=value1,key2=value2"

const testSystemRoutesFlag = true
const testDynamicRoutesFlag = true

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
		systemRoutesFlag:       strconv.FormatBool(testSystemRoutesFlag),
		dynamicRoutesFlag:      strconv.FormatBool(testDynamicRoutesFlag),
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
		Name:           utils.Ptr(testRoutingTableName),
		Description:    utils.Ptr(testRoutingTableDescription),
		SystemRoutes:   testSystemRoutesFlag,
		DynamicRoutes:  testDynamicRoutesFlag,
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

func fixtureRequest(mods ...func(request *iaas.ApiUpdateRoutingTableOfAreaRequest)) iaas.ApiUpdateRoutingTableOfAreaRequest {
	req := testClient.UpdateRoutingTableOfArea(
		testCtx,
		testOrgId,
		testNetworkAreaId,
		testRegion,
		testRoutingTableId,
	)

	payload := iaas.UpdateRoutingTableOfAreaPayload{
		Labels:        utils.ConvertStringMapToInterfaceMap(testLabels),
		Name:          utils.Ptr(testRoutingTableName),
		Description:   utils.Ptr(testRoutingTableDescription),
		DynamicRoutes: utils.Ptr(true),
		SystemRoutes:  utils.Ptr(true),
	}

	req = req.UpdateRoutingTableOfAreaPayload(payload)

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
			description: "base",
			flagValues:  fixtureFlagValues(),
			argValues:   fixtureArgValues(),
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = testRoutingTableId
			}),
		},
		{
			description: "dynamic routes disabled",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[dynamicRoutesFlag] = "false"
			}),
			argValues: fixtureArgValues(),
			isValid:   true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DynamicRoutes = false
				model.RoutingTableId = testRoutingTableId
			}),
		},
		{
			description: "system routes disabled",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[systemRoutesFlag] = "false"
			}),
			argValues: fixtureArgValues(),
			isValid:   true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SystemRoutes = false
				model.RoutingTableId = testRoutingTableId
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

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiUpdateRoutingTableOfAreaRequest
	}{
		{
			description: "base",
			model: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = testRoutingTableId
			}),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "labels missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = testRoutingTableId
				model.Labels = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateRoutingTableOfAreaRequest) {
				*request = (*request).UpdateRoutingTableOfAreaPayload(iaas.UpdateRoutingTableOfAreaPayload{
					Labels:        nil,
					Name:          utils.Ptr(testRoutingTableName),
					Description:   utils.Ptr(testRoutingTableDescription),
					DynamicRoutes: utils.Ptr(true),
					SystemRoutes:  utils.Ptr(true),
				})
			}),
		},
		{
			description: "name missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = testRoutingTableId
				model.Name = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateRoutingTableOfAreaRequest) {
				*request = (*request).UpdateRoutingTableOfAreaPayload(iaas.UpdateRoutingTableOfAreaPayload{
					Labels:        utils.ConvertStringMapToInterfaceMap(testLabels),
					Name:          nil,
					Description:   utils.Ptr(testRoutingTableDescription),
					DynamicRoutes: utils.Ptr(true),
					SystemRoutes:  utils.Ptr(true),
				})
			}),
		},
		{
			description: "description missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = testRoutingTableId
				model.Description = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateRoutingTableOfAreaRequest) {
				*request = (*request).UpdateRoutingTableOfAreaPayload(iaas.UpdateRoutingTableOfAreaPayload{
					Labels:        utils.ConvertStringMapToInterfaceMap(testLabels),
					Name:          utils.Ptr(testRoutingTableName),
					Description:   nil,
					DynamicRoutes: utils.Ptr(true),
					SystemRoutes:  utils.Ptr(true),
				})
			}),
		},
		{
			description: "dynamic routes disabled",
			model: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = testRoutingTableId
				model.DynamicRoutes = false
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateRoutingTableOfAreaRequest) {
				*request = (*request).UpdateRoutingTableOfAreaPayload(iaas.UpdateRoutingTableOfAreaPayload{
					Labels:        utils.ConvertStringMapToInterfaceMap(testLabels),
					Name:          utils.Ptr(testRoutingTableName),
					Description:   utils.Ptr(testRoutingTableDescription),
					DynamicRoutes: utils.Ptr(false),
					SystemRoutes:  utils.Ptr(true),
				})
			}),
		},
		{
			description: "system routes disabled",
			model: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableId = testRoutingTableId
				model.DynamicRoutes = false
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateRoutingTableOfAreaRequest) {
				*request = (*request).UpdateRoutingTableOfAreaPayload(iaas.UpdateRoutingTableOfAreaPayload{
					Labels:        utils.ConvertStringMapToInterfaceMap(testLabels),
					Name:          utils.Ptr(testRoutingTableName),
					Description:   utils.Ptr(testRoutingTableDescription),
					SystemRoutes:  utils.Ptr(true),
					DynamicRoutes: utils.Ptr(false),
				})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := buildRequest(testCtx, tt.model, testClient)

			if diff := cmp.Diff(req, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			); diff != "" {
				t.Errorf("buildRequest() mismatch (-got +want):\n%s", diff)
			}
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
			outputFormat: print.PrettyOutputFormat,
			routingTable: nil,
			wantErr:      true,
		},
		{
			name:         "empty routing-table",
			outputFormat: print.PrettyOutputFormat,
			routingTable: &iaas.RoutingTable{},
			wantErr:      true,
		},
		{
			name:         "pretty output routing-table",
			outputFormat: print.PrettyOutputFormat,
			routingTable: &dummyRoutingTable,
			wantErr:      false,
		},
		{
			name:         "json output routing-table",
			outputFormat: print.JSONOutputFormat,
			routingTable: &dummyRoutingTable,
			wantErr:      false,
		},
		{
			name:         "yaml output routing-table",
			outputFormat: print.YAMLOutputFormat,
			routingTable: &dummyRoutingTable,
			wantErr:      false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.outputFormat, "network-area-id", tt.routingTable); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
