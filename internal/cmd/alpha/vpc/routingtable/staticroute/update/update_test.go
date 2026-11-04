package update

import (
	"context"
	"maps"
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
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &iaas.APIClient{DefaultAPI: &iaas.DefaultAPIService{}}

	testProjectId      = uuid.NewString()
	testVpcId          = uuid.NewString()
	testRoutingTableId = uuid.NewString()
	testRouteId        = uuid.NewString()

	testLabels = map[string]any{
		"key1": "value1",
		"key2": "value2",
	}
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testRouteId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		vpcIdFlag:                 testVpcId,
		routingTableIdFlag:        testRoutingTableId,
		labelFlag:                 testLabelSelector,
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
		RouteId:        testRouteId,
		RoutingTableId: testRoutingTableId,
		VpcId:          testVpcId,
		Labels:         maps.Clone(testLabels),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(req *iaas.ApiUpdateVPCStaticRouteRequest)) iaas.ApiUpdateVPCStaticRouteRequest {
	req := testClient.DefaultAPI.UpdateVPCStaticRoute(testCtx, testProjectId, testVpcId, testRegion, testRoutingTableId, testRouteId)
	payload := iaas.UpdateVPCStaticRoutePayload{
		Labels: testLabels,
	}
	req = req.UpdateVPCStaticRoutePayload(payload)

	for _, mod := range mods {
		mod(&req)
	}

	return req
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
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
			description:   "arg value invalid 1",
			argValues:     []string{""},
			flagValues:    fixtureFlagValues(),
			isValid:       false,
			expectedModel: fixtureInputModel(),
		},
		{
			description:   "arg value invalid 2",
			argValues:     []string{"invalid-uuid"},
			flagValues:    fixtureFlagValues(),
			isValid:       false,
			expectedModel: fixtureInputModel(),
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
			description: "vpc-id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, vpcIdFlag)
			}),
			isValid: false,
		},
		{
			description: "labels are missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid label format",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelFlag] = "invalid-label"
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
		expectedRequest iaas.ApiUpdateVPCStaticRouteRequest
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
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateVPCStaticRouteRequest) {
				*request = request.UpdateVPCStaticRoutePayload(iaas.UpdateVPCStaticRoutePayload{
					Labels: nil,
				})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			gotReq := buildRequest(testCtx, tt.model, testClient.DefaultAPI)

			if diff := cmp.Diff(tt.expectedRequest, gotReq,
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
		Labels:    testLabels,
		CreatedAt: new(time.Now()),
		UpdatedAt: new(time.Now()),
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
			wantErr:      true,
		},
		{
			name:         "pretty output with one route",
			outputFormat: print.PrettyOutputFormat,
			route:        new(fixtureRoute),
			wantErr:      false,
		},
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			route:        new(fixtureRoute),
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			route:        new(fixtureRoute),
			wantErr:      false,
		},
	}

	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, "", "", tt.route); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
