package delete

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
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
		args          []string
		flagValues    map[string]string
		isValid       bool
		expectedRoute *inputModel
	}{
		{
			description: "valid input",
			args:        []string{testRouteId},
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expectedRoute: fixtureInputModel(func(m *inputModel) {
				m.RouteID = &testRouteId
			}),
		},
		{
			description: "missing route id arg",
			args:        []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "missing organization-id flag",
			args:        []string{testRouteId},
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, "organization-id")
			}),
			isValid: false,
		},
		{
			description: "missing network-area-id flag",
			args:        []string{testRouteId},
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, "network-area-id")
			}),
			isValid: false,
		},
		{
			description: "missing routing-table-id flag",
			args:        []string{testRouteId},
			flagValues: fixtureFlagValues(func(m map[string]string) {
				delete(m, "routing-table-id")
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.args)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedRoute)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
