package create

import (
	"context"
	"maps"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &iaas.APIClient{DefaultAPI: &iaas.DefaultAPIService{}}

	testProjectId      = uuid.NewString()
	testVpcId          = uuid.NewString()
	testRoutingTableId = uuid.NewString()
	testLabels         = map[string]any{
		"key1": "value1",
		"key2": "value2",
	}
)

const (
	testRegion           = "eu01"
	testDestinationType  = destTypeCIDRv4
	testDestinationValue = "1.1.1.0/24"
	testNextHopType      = nextHopTypeIPv4
	testNextHopValue     = "1.1.1.1"
	testLabelSelector    = "key1=value1,key2=value2"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag:  testProjectId,
		globalflags.RegionFlag:     testRegion,
		labelFlag:                  testLabelSelector,
		vpcIdFlag:                  testVpcId,
		routingTableIdFlag:         testRoutingTableId,
		destinationTypeFlag.Name(): testDestinationType,
		destinationValueFlag:       testDestinationValue,
		nextHopTypeFlag.Name():     testNextHopType,
		nextHopValueFlag:           testNextHopValue,
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
		VpcId:            testVpcId,
		RoutingTableId:   testRoutingTableId,
		DestinationType:  testDestinationType,
		DestinationValue: testDestinationValue,
		NextHopType:      testNextHopType,
		NextHopValue:     new(testNextHopValue),
		Labels:           maps.Clone(testLabels),
	}
	for _, mod := range mods {
		mod(model)
	}

	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiAddVPCStaticRouteRequest)) iaas.ApiAddVPCStaticRouteRequest {
	request := testClient.DefaultAPI.AddVPCStaticRoute(testCtx, testProjectId, testVpcId, testRegion, testRoutingTableId)

	request = request.AddVPCStaticRoutePayload(fixturePayload())

	for _, mod := range mods {
		mod(&request)
	}

	return request
}

func fixturePayload(mods ...func(payload *iaas.AddVPCStaticRoutePayload)) iaas.AddVPCStaticRoutePayload {
	payload := iaas.AddVPCStaticRoutePayload{
		Destination: iaas.AddVPCStaticRoutePayloadDestination{
			DestinationCIDRv4: &iaas.DestinationCIDRv4{
				Type:  testDestinationType,
				Value: testDestinationValue,
			},
		},
		Nexthop: iaas.AddVPCStaticRoutePayloadNexthop{
			NexthopIPv4: &iaas.NexthopIPv4{
				Type:  testNextHopType,
				Value: testNextHopValue,
			},
		},
		Labels: testLabels,
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
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "routing-table ID missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, routingTableIdFlag)
			}),
			isValid: false,
		},
		{
			description: "destination value missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, destinationValueFlag)
			}),
			isValid: false,
		},
		{
			description: "destination type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, destinationTypeFlag.Name())
			}),
			isValid: false,
		},
		{
			description: "next hop type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nextHopTypeFlag.Name())
			}),
			isValid: false,
		},
		{
			description: "next hop value missing when nexthop type is ipv4",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag.Name()] = testNextHopType
				delete(flagValues, nextHopValueFlag)
			}),
			isValid: false,
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "vpc ID missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, vpcIdFlag)
			}),
			isValid: false,
		},
		{
			description: "vpc ID invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[vpcIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "routing table ID invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[routingTableIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "invalid destination type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationTypeFlag.Name()] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "invalid next hop type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag.Name()] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "valid input - internet next hop without next hop value",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag.Name()] = nextHopTypeInternet
				delete(flagValues, nextHopValueFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = nextHopTypeInternet
				model.NextHopValue = nil
			}),
		},
		{
			description: "invalid input - internet next hop with next hop value",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag.Name()] = nextHopTypeInternet
				flagValues[nextHopValueFlag] = "some-value"
			}),
			isValid: false,
		},
		{
			description: "valid input - blackhole next hop without next hop value",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag.Name()] = nextHopTypeBlackhole
				delete(flagValues, nextHopValueFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = nextHopTypeBlackhole
				model.NextHopValue = nil
			}),
		},
		{
			description: "invalid input - blackhole next hop with next hop value",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag.Name()] = nextHopTypeBlackhole
				flagValues[nextHopValueFlag] = "some-value"
			}),
			isValid: false,
		},
		{
			description: "valid input - CIDRv6 destination and IPv6 nexthop",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationTypeFlag.Name()] = destTypeCIDRv6
				flagValues[destinationValueFlag] = "2001:db8::/32"
				flagValues[nextHopTypeFlag.Name()] = nextHopTypeIPv6
				flagValues[nextHopValueFlag] = "2001:db8::1"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DestinationType = destTypeCIDRv6
				model.DestinationValue = "2001:db8::/32"
				model.NextHopType = nextHopTypeIPv6
				model.NextHopValue = utils.Ptr("2001:db8::1")
			}),
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
		expectedRequest iaas.ApiAddVPCStaticRouteRequest
		expectErr       bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "cidrv6 and ipv6",
			model: fixtureInputModel(func(model *inputModel) {
				model.DestinationType = destTypeCIDRv6
				model.DestinationValue = "2001:db8::/32"
				model.NextHopType = nextHopTypeIPv6
				model.NextHopValue = utils.Ptr("2001:db8::1")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddVPCStaticRouteRequest) {
				payload := fixturePayload(func(payload *iaas.AddVPCStaticRoutePayload) {
					payload.Destination = iaas.AddVPCStaticRoutePayloadDestination{
						DestinationCIDRv6: &iaas.DestinationCIDRv6{
							Type:  destTypeCIDRv6,
							Value: "2001:db8::/32",
						},
					}
					payload.Nexthop = iaas.AddVPCStaticRoutePayloadNexthop{
						NexthopIPv6: &iaas.NexthopIPv6{
							Type:  nextHopTypeIPv6,
							Value: "2001:db8::1",
						},
					}
				})

				*request = request.AddVPCStaticRoutePayload(payload)
			}),
		},
		{
			description: "nexthop internet",
			model: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = nextHopTypeInternet
				model.NextHopValue = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddVPCStaticRouteRequest) {
				payload := fixturePayload(func(payload *iaas.AddVPCStaticRoutePayload) {
					payload.Nexthop = iaas.AddVPCStaticRoutePayloadNexthop{
						NexthopInternet: &iaas.NexthopInternet{
							Type: nextHopTypeInternet,
						},
					}
				})

				*request = request.AddVPCStaticRoutePayload(payload)
			}),
		},
		{
			description: "nexthop blackhole",
			model: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = nextHopTypeBlackhole
				model.NextHopValue = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddVPCStaticRouteRequest) {
				payload := fixturePayload(func(payload *iaas.AddVPCStaticRoutePayload) {
					payload.Nexthop = iaas.AddVPCStaticRoutePayloadNexthop{
						NexthopBlackhole: &iaas.NexthopBlackhole{
							Type: nextHopTypeBlackhole,
						},
					}
				})

				*request = request.AddVPCStaticRoutePayload(payload)
			}),
		},
		{
			description: "invalid destination type returns error",
			model: fixtureInputModel(func(model *inputModel) {
				model.DestinationType = "unsupported"
			}),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			gotReq, err := buildRequest(testCtx, tt.model, testClient.DefaultAPI)
			if (err != nil) != tt.expectErr {
				t.Fatalf("buildRequest() error = %v, expectErr %v", err, tt.expectErr)
			}

			if tt.expectErr {
				return
			}

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
	tests := []struct {
		name         string
		outputFormat string
		route        *iaas.Route
		wantErr      bool
	}{
		{
			name:    "nil route should return error",
			route:   nil,
			wantErr: true,
		},
		{
			name:    "empty route",
			route:   &iaas.Route{},
			wantErr: false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, tt.route); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
