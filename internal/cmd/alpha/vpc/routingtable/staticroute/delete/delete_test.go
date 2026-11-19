package delete

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

const (
	testRegion = "eu01"
)

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &iaas.APIClient{DefaultAPI: &iaas.DefaultAPIService{}}

	testProjectId      = uuid.NewString()
	testVpcId          = uuid.NewString()
	testRoutingTableId = uuid.NewString()
	testRouteId        = uuid.NewString()
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

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		vpcIdFlag:                 testVpcId,
		routingTableIdFlag:        testRoutingTableId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}

	return flagValues
}

func fixtureInputModel(mods ...func(*inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		RouteId:        testRouteId,
		RoutingTableId: testRoutingTableId,
		VpcId:          testVpcId,
	}
	for _, mod := range mods {
		mod(model)
	}

	return model
}

func fixtureRequest(mods ...func(*iaas.ApiDeleteVPCStaticRouteRequest)) iaas.ApiDeleteVPCStaticRouteRequest {
	req := testClient.DefaultAPI.DeleteVPCStaticRoute(testCtx, testProjectId, testVpcId, testRegion, testRoutingTableId, testRouteId)
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
			description: "missing route id arg",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "missing vpc-id flag",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, vpcIdFlag)
			}),
			isValid: false,
		},
		{
			description: "missing routing-table-id flag",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, routingTableIdFlag)
			}),
			isValid: false,
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
			description: "invalid vpc-id flag",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flags map[string]string) {
				flags[vpcIdFlag] = "invalid-vpc"
			}),
			isValid: false,
		},
		{
			description: "invalid routing-table-id flag",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flags map[string]string) {
				flags[routingTableIdFlag] = "invalid-table"
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
		expectedRequest iaas.ApiDeleteVPCStaticRouteRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
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
