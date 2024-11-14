package update

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRouteId = uuid.NewString()

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
		organizationIdFlag: testOrgId,
		networkAreaIdFlag:  testNetworkAreaId,
		labelFlag:          "value=key",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixturePayload(mods ...func(payload *iaas.UpdateNetworkAreaRoutePayload)) iaas.UpdateNetworkAreaRoutePayload {
	payload := iaas.UpdateNetworkAreaRoutePayload{
		Labels: &map[string]interface{}{
			"value": "key",
		},
	}

	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixturePayloadAsStringMap() map[string]string {
	payload := fixturePayload()
	labelsMap := make(map[string]string)
	for k, v := range *payload.Labels {
		if value, ok := v.(string); ok {
			labelsMap[k] = value
		}
	}
	return labelsMap
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	payload := fixturePayloadAsStringMap()
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		OrganizationId: utils.Ptr(testOrgId),
		NetworkAreaId:  utils.Ptr(testNetworkAreaId),
		RouteId:        testRouteId,
		Labels:         utils.Ptr(payload),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiUpdateNetworkAreaRouteRequest)) iaas.ApiUpdateNetworkAreaRouteRequest {
	request := testClient.UpdateNetworkAreaRoute(testCtx, testOrgId, testNetworkAreaId, testRouteId)
	request = request.UpdateNetworkAreaRoutePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "org id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "org id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "org area id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "network area id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "network area id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "network area id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "route id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, routeIdArg)
			}),
			isValid: false,
		},
		{
			description: "route id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[routeIdArg] = ""
			}),
			isValid: false,
		},
		{
			description: "route id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[routeIdArg] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "labels missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
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

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiUpdateNetworkAreaRouteRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}