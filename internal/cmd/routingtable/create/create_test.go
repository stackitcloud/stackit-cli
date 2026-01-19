package create

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

const testRoutingTableName = "test"
const testRoutingTableDescription = "test"

const testSystemRoutesFlag = true
const testDynamicRoutesFlag = true

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
		Name:           testRoutingTableName,
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

func fixtureRequest(mods ...func(request *iaas.ApiAddRoutingTableToAreaRequest)) iaas.ApiAddRoutingTableToAreaRequest {
	request := testClient.AddRoutingTableToArea(testCtx, testOrgId, testNetworkAreaId, testRegion)
	request = request.AddRoutingTableToAreaPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.AddRoutingTableToAreaPayload)) iaas.AddRoutingTableToAreaPayload {
	payload := iaas.AddRoutingTableToAreaPayload{
		Description:   utils.Ptr(testRoutingTableDescription),
		Name:          utils.Ptr(testRoutingTableName),
		Labels:        utils.ConvertStringMapToInterfaceMap(testLabels),
		SystemRoutes:  utils.Ptr(true),
		DynamicRoutes: utils.Ptr(true),
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
			description:   "valid input",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "dynamic routes disabled",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[dynamicRoutesFlag] = "false"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DynamicRoutes = false
			}),
		},
		{
			description: "system routes disabled",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[systemRoutesFlag] = "false"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SystemRoutes = false
			}),
		},
		{
			description: "missing organization ID",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid organization ID - empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "invalid organization ID - format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "missing network area ID",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid network area ID - empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "invalid network area ID - format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "missing name",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "missing labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
		},
		{
			description: "missing description",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, descriptionFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
		},
		{
			description: "no flags provided",
			flagValues:  map[string]string{},
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
		expectedRequest iaas.ApiAddRoutingTableToAreaRequest
	}{
		{
			description:     "valid input",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "labels missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutingTableToAreaRequest) {
				*request = (*request).AddRoutingTableToAreaPayload(
					fixturePayload(func(payload *iaas.AddRoutingTableToAreaPayload) {
						payload.Labels = nil
					}),
				)
			}),
		},
		{
			description: "system routes disabled",
			model: fixtureInputModel(func(model *inputModel) {
				model.SystemRoutes = false
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutingTableToAreaRequest) {
				*request = (*request).AddRoutingTableToAreaPayload(
					fixturePayload(func(payload *iaas.AddRoutingTableToAreaPayload) {
						payload.SystemRoutes = utils.Ptr(false)
					}),
				)
			}),
		},
		{
			description: "dynamic routes disabled",
			model: fixtureInputModel(func(model *inputModel) {
				model.DynamicRoutes = false
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutingTableToAreaRequest) {
				*request = (*request).AddRoutingTableToAreaPayload(
					fixturePayload(func(payload *iaas.AddRoutingTableToAreaPayload) {
						payload.DynamicRoutes = utils.Ptr(false)
					}),
				)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				t.Fatalf("buildRequest returned error: %v", err)
			}

			if diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx)); diff != "" {
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
			outputFormat: "",
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
			if err := outputResult(p, tt.outputFormat, tt.routingTable); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
