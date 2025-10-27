package delete

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

var (
	testOrgId          = uuid.NewString()
	testNetworkAreaId  = uuid.NewString()
	testRoutingTableId = uuid.NewString()
	testRouteId        = uuid.NewString()
)

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	flagValues := map[string]string{
		organizationIdFlag: testOrgId,
		networkAreaIdFlag:  testNetworkAreaId,
		routingTableIdFlag: testRoutingTableId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(*inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.InfoVerbosity,
		},
		OrganizationId: &testOrgId,
		NetworkAreaId:  &testNetworkAreaId,
		RoutingTableId: &testRoutingTableId,
		RouteID:        &testRouteId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			description: "valid input",
			argValues:   []string{testRouteId},
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expectedModel: fixtureInputModel(func(m *inputModel) {
				m.RouteID = &testRouteId
			}),
		},
		{
			description: "missing route id arg",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "missing organization-id flag",
			argValues:   []string{testRouteId},
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, "organization-id")
			}),
			isValid: false,
		},
		{
			description: "missing network-area-id flag",
			argValues:   []string{testRouteId},
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, "network-area-id")
			}),
			isValid: false,
		},
		{
			description: "missing routing-table-id flag",
			argValues:   []string{testRouteId},
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, "routing-table-id")
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
			description: "invalid organization-id flag",
			argValues:   []string{testRouteId},
			flagValues:  map[string]string{"organization-id": "invalid-org"},
			isValid:     false,
		},
		{
			description: "invalid network-area-id flag",
			argValues:   []string{testRouteId},
			flagValues:  map[string]string{"network-area-id": "invalid-area"},
			isValid:     false,
		},
		{
			description: "invalid routing-table-id flag",
			argValues:   []string{testRouteId},
			flagValues:  map[string]string{"routing-table-id": "invalid-table"},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}
