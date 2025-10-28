package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	testRegion            = "eu01"
	testDestinationCIDRv4 = "1.1.1.0/24"
	testNexthopIPv4       = "1.1.1.1"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,

		organizationIdFlag: testOrgId,
		networkAreaIdFlag:  testNetworkAreaId,
		destinationFlag:    testDestinationCIDRv4,
		nexthopIPv4Flag:    testNexthopIPv4,
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
		DestinationV4:  utils.Ptr(testDestinationCIDRv4),
		NexthopV4:      utils.Ptr(testNexthopIPv4),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateNetworkAreaRouteRequest)) iaas.ApiCreateNetworkAreaRouteRequest {
	request := testClient.CreateNetworkAreaRoute(testCtx, testOrgId, testNetworkAreaId, testRegion)
	request = request.CreateNetworkAreaRoutePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateNetworkAreaRoutePayload)) iaas.CreateNetworkAreaRoutePayload {
	payload := iaas.CreateNetworkAreaRoutePayload{
		Items: &[]iaas.Route{
			{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(destinationCIDRv4Type),
						Value: utils.Ptr(testDestinationCIDRv4),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr(nexthopIPv4Type),
						Value: utils.Ptr(testNexthopIPv4),
					},
				},
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
		argValues     []string
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
				delete(flagValues, nexthopIPv4Flag)
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
			description: "destination missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, destinationFlag)
			}),
			isValid: false,
		},
		{
			description: "destinationFlag invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "destinationFlag invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationFlag] = "invalid-destinationFlag"
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
		{
			description: "conflicting destination and prefix set",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[prefixFlag] = testDestinationCIDRv4
			}),
			isValid: false,
		},
		{
			description: "conflicting nexthop and nexthop-ipv4 set",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nexthopFlag] = testNexthopIPv4
			}),
			isValid: false,
		},
		{
			description: "conflicting nexthop and nexthop-ipv4 set",
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
				*request = (*request).CreateNetworkAreaRoutePayload(fixturePayload(func(payload *iaas.CreateNetworkAreaRoutePayload) {
					(*payload.Items)[0].Labels = utils.Ptr(map[string]interface{}{"key": "value"})
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

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat     string
		networkAreaLabel string
		route            iaas.Route
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "empty route",
			args: args{
				route: iaas.Route{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.networkAreaLabel, tt.args.route); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
