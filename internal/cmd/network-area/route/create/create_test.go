package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		organizationIdFlag: testOrgId,
		networkAreaIdFlag:  testNetworkAreaId,
		prefixFlag:         "1.1.1.0/24",
		nexthopFlag:        "1.1.1.1",
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
		},
		OrganizationId: utils.Ptr(testOrgId),
		NetworkAreaId:  utils.Ptr(testNetworkAreaId),
		Prefix:         utils.Ptr("1.1.1.0/24"),
		Nexthop:        utils.Ptr("1.1.1.1"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateNetworkAreaRouteRequest)) iaas.ApiCreateNetworkAreaRouteRequest {
	request := testClient.CreateNetworkAreaRoute(testCtx, testOrgId, testNetworkAreaId)
	request = request.CreateNetworkAreaRoutePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateNetworkAreaRoutePayload)) iaas.CreateNetworkAreaRoutePayload {
	payload := iaas.CreateNetworkAreaRoutePayload{
		Ipv4: &[]iaas.Route{
			{
				Prefix:  utils.Ptr("1.1.1.0/24"),
				Nexthop: utils.Ptr("1.1.1.1"),
			},
		},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		aclValues     []string
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
			description: "next hop missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nexthopFlag)
			}),
			isValid: false,
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
			description: "prefix missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, prefixFlag)
			}),
			isValid: false,
		},
		{
			description: "prefix invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[prefixFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "prefix invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[prefixFlag] = "invalid-prefix"
			}),
			isValid: false,
		},
		{
			description: "optional labels is provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelFlag] = "key=value"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = utils.Ptr(map[string]string{"key": "value"})
			}),
			isValid: true,
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

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd)
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
		expectedRequest iaas.ApiCreateNetworkAreaRouteRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "optional labels provided",
			model: fixtureInputModel(func(model *inputModel) {
				model.Labels = utils.Ptr(map[string]string{"key": "value"})
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateNetworkAreaRouteRequest) {
				*request = request.CreateNetworkAreaRoutePayload(fixturePayload(func(payload *iaas.CreateNetworkAreaRoutePayload) {
					(*payload.Ipv4)[0].Labels = utils.Ptr(map[string]interface{}{"key": "value"})
				}))
			}),
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
