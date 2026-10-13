package list

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

const (
	testRegion        = "eu01"
	testLabelSelector = "key1=value1,key2=value2"
)

var (
	testCtx            = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient         = &iaas.APIClient{DefaultAPI: &iaas.DefaultAPIService{}}
	testProjectId      = uuid.NewString()
	testVpcId          = uuid.NewString()
	testRoutingTableId = uuid.NewString()

	testLimit int64 = 10
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		vpcIdFlag:                 testVpcId,
		routingTableIdFlag:        testRoutingTableId,
		labelSelectorFlag:         testLabelSelector,
		limitFlag:                 strconv.Itoa(int(testLimit)),
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
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		RoutingTableId: testRoutingTableId,
		VpcId:          testVpcId,
		LabelSelector:  new(testLabelSelector),
		Limit:          new(testLimit),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiListVPCStaticRoutesRequest)) iaas.ApiListVPCStaticRoutesRequest {
	request := testClient.DefaultAPI.ListVPCStaticRoutes(testCtx, testProjectId, testVpcId, testRegion, testRoutingTableId)
	request = request.LabelSelector(testLabelSelector)
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			description: "vpc-id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, vpcIdFlag)
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

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiListVPCStaticRoutesRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "label selector nil",
			model: fixtureInputModel(func(m *inputModel) {
				m.LabelSelector = nil
			}),
			expectedRequest: testClient.DefaultAPI.ListVPCStaticRoutes(testCtx, testProjectId, testVpcId, testRegion, testRoutingTableId),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			gotReq := buildRequest(testCtx, tt.model, testClient.DefaultAPI)

			if diff := cmp.Diff(
				tt.expectedRequest,
				gotReq,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, iaas.DefaultAPIService{}),
			); diff != "" {
				t.Errorf("buildRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	fixtureRoute := iaas.Route{
		Id: new("route-foo"),
		Destination: iaas.RouteDestination{
			DestinationCIDRv4: &iaas.DestinationCIDRv4{
				Type:  "cidrv4",
				Value: "10.0.0.0/24",
			},
		},
		Nexthop: iaas.RouteNexthop{
			NexthopIPv4: &iaas.NexthopIPv4{
				Type:  "ipv4",
				Value: "10.0.0.1",
			},
		},
		Labels: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
		CreatedAt: new(time.Now()),
		UpdatedAt: new(time.Now()),
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
			name:         "empty routes list yaml output",
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
			routes:       []iaas.Route{fixtureRoute},
			wantErr:      false,
		},
		{
			name:         "pretty output with multiple routes",
			outputFormat: print.PrettyOutputFormat,
			routes:       []iaas.Route{fixtureRoute, fixtureRoute, fixtureRoute},
			wantErr:      false,
		},
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			routes:       []iaas.Route{fixtureRoute},
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			routes:       []iaas.Route{fixtureRoute},
			wantErr:      false,
		},
	}

	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, tt.routes, "dummy-vpc-id", "dummy-route-table-id"); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
