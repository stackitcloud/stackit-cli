package create

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

const testRegion = "eu01"

var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRoutingTableId = uuid.NewString()

const testDestinationTypeFlag = "cidrv4"
const testDestinationValueFlag = "1.1.1.0/24"
const testNextHopTypeFlag = "ipv4"
const testNextHopValueFlag = "1.1.1.1"
const testLabelSelectorFlag = "key1=value1,key2=value2"

var testLabels = &map[string]string{
	"key1": "value1",
	"key2": "value2",
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,
		labelFlag:              testLabelSelectorFlag,
		organizationIdFlag:     testOrgId,
		networkAreaIdFlag:      testNetworkAreaId,
		routingTableIdFlag:     testRoutingTableId,
		destinationTypeFlag:    testDestinationTypeFlag,
		destinationValueFlag:   testDestinationValueFlag,
		nextHopTypeFlag:        testNextHopTypeFlag,
		nextHopValueFlag:       testNextHopValueFlag,
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
		OrganizationId:   testOrgId,
		NetworkAreaId:    testNetworkAreaId,
		RoutingTableId:   testRoutingTableId,
		DestinationType:  utils.Ptr(testDestinationTypeFlag),
		DestinationValue: utils.Ptr(testDestinationValueFlag),
		NextHopType:      utils.Ptr(testNextHopTypeFlag),
		NextHopValue:     utils.Ptr(testNextHopValueFlag),
		Labels:           testLabels,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiAddRoutesToRoutingTableRequest)) iaas.ApiAddRoutesToRoutingTableRequest {
	request := testClient.AddRoutesToRoutingTable(testCtx, testOrgId, testNetworkAreaId, testRegion, testRoutingTableId)
	request = request.AddRoutesToRoutingTablePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.AddRoutesToRoutingTablePayload)) iaas.AddRoutesToRoutingTablePayload {
	payload := iaas.AddRoutesToRoutingTablePayload{
		Items: &[]iaas.Route{
			{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(testDestinationTypeFlag),
						Value: utils.Ptr(testDestinationValueFlag),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr(testNextHopTypeFlag),
						Value: utils.Ptr(testNextHopValueFlag),
					},
				},
				Labels: utils.ConvertStringMapToInterfaceMap(testLabels),
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
				delete(flagValues, destinationTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "next hop type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nextHopTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "next hop value missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
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
			description: "organization ID missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "organization ID invalid - empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "organization ID invalid - format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "network area ID missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "network area ID invalid - empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "network area ID invalid - format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "invalid destination type enum",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationTypeFlag] = "ipv4"
			}),
			isValid: false,
		},
		{
			description: "destination value not IPv4 CIDR",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationValueFlag] = "0.0.0.0"
			}),
			isValid: false,
		},
		{
			description: "destination value not IPv6 CIDR",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationTypeFlag] = "cidrv6"
				flagValues[destinationValueFlag] = "2001:db8::"
			}),
			isValid: false,
		},
		{
			description: "destination value is IPv6 CIDR",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationTypeFlag] = "cidrv6"
				flagValues[destinationValueFlag] = "2001:db8::/32"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DestinationType = utils.Ptr("cidrv6")
				model.DestinationValue = utils.Ptr("2001:db8::/32")
			}),
		},
		{
			description: "invalid next hop type enum",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "cidrv4"
			}),
			isValid: false,
		},
		{
			description: "next hop type is internet and next hop value is provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "internet"
				flagValues[nextHopValueFlag] = "1.1.1.1" // should not be allowed
			}),
			isValid: false,
		},
		{
			description: "next hop type is blackhole and next hop value is provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "blackhole"
				flagValues[nextHopValueFlag] = "1.1.1.1"
			}),
			isValid: false,
		},
		{
			description: "next hop type is internet and next hop value is not provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "internet"
				delete(flagValues, nextHopValueFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = utils.Ptr("internet")
				model.NextHopValue = nil
			}),
			isValid: true,
		},
		{
			description: "next hop type is blackhole and next hop value is not provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "blackhole"
				delete(flagValues, nextHopValueFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = utils.Ptr("blackhole")
				model.NextHopValue = nil
			}),
			isValid: true,
		},
		{
			description: "next hop type is IPv4 and next hop value is missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "ipv4"
				delete(flagValues, nextHopValueFlag)
			}),
			isValid: false,
		},
		{
			description: "next hop type is IPv6 and next hop value is missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "ipv6"
				delete(flagValues, nextHopValueFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid next hop type provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "invalid-type"
			}),
			isValid: false,
		},
		{
			description: "optional labels are provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelFlag] = "key=value"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = utils.Ptr(map[string]string{"key": "value"})
			}),
			isValid: true,
		},
		{
			description:   "argument value missing",
			argValues:     []string{""},
			flagValues:    fixtureFlagValues(),
			isValid:       false,
			expectedModel: fixtureInputModel(),
		},
		{
			description:   "argument value wrong",
			argValues:     []string{"foo-bar"},
			flagValues:    fixtureFlagValues(),
			isValid:       false,
			expectedModel: fixtureInputModel(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildNextHop(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expected    *iaas.RouteNexthop
	}{
		{
			description: "IPv4 next hop",
			model: fixtureInputModel(func(m *inputModel) {
				m.NextHopType = utils.Ptr("ipv4")
				m.NextHopValue = utils.Ptr("1.1.1.1")
			}),
			expected: &iaas.RouteNexthop{
				NexthopIPv4: &iaas.NexthopIPv4{
					Type:  utils.Ptr("ipv4"),
					Value: utils.Ptr("1.1.1.1"),
				},
			},
		},
		{
			description: "IPv6 next hop",
			model: fixtureInputModel(func(m *inputModel) {
				m.NextHopType = utils.Ptr("ipv6")
				m.NextHopValue = utils.Ptr("::1")
			}),
			expected: &iaas.RouteNexthop{
				NexthopIPv6: &iaas.NexthopIPv6{
					Type:  utils.Ptr("ipv6"),
					Value: utils.Ptr("::1"),
				},
			},
		},
		{
			description: "Internet next hop",
			model: fixtureInputModel(func(m *inputModel) {
				m.NextHopType = utils.Ptr("internet")
				m.NextHopValue = nil
			}),
			expected: &iaas.RouteNexthop{
				NexthopInternet: &iaas.NexthopInternet{
					Type: utils.Ptr("internet"),
				},
			},
		},
		{
			description: "Blackhole next hop",
			model: fixtureInputModel(func(m *inputModel) {
				m.NextHopType = utils.Ptr("blackhole")
				m.NextHopValue = nil
			}),
			expected: &iaas.RouteNexthop{
				NexthopBlackhole: &iaas.NexthopBlackhole{
					Type: utils.Ptr("blackhole"),
				},
			},
		},
		{
			description: "Unsupported next hop type",
			model: fixtureInputModel(func(m *inputModel) {
				m.NextHopType = utils.Ptr("unsupported")
			}),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := buildNextHop(tt.model)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("buildNextHop() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBuildDestination(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expected    *iaas.RouteDestination
	}{
		{
			description: "CIDRv4 destination",
			model: fixtureInputModel(func(m *inputModel) {
				m.DestinationType = utils.Ptr("cidrv4")
				m.DestinationValue = utils.Ptr("192.168.1.0/24")
			}),
			expected: &iaas.RouteDestination{
				DestinationCIDRv4: &iaas.DestinationCIDRv4{
					Type:  utils.Ptr("cidrv4"),
					Value: utils.Ptr("192.168.1.0/24"),
				},
			},
		},
		{
			description: "CIDRv6 destination",
			model: fixtureInputModel(func(m *inputModel) {
				m.DestinationType = utils.Ptr("cidrv6")
				m.DestinationValue = utils.Ptr("2001:db8::/32")
			}),
			expected: &iaas.RouteDestination{
				DestinationCIDRv6: &iaas.DestinationCIDRv6{
					Type:  utils.Ptr("cidrv6"),
					Value: utils.Ptr("2001:db8::/32"),
				},
			},
		},
		{
			description: "unsupported destination type",
			model: fixtureInputModel(func(m *inputModel) {
				m.DestinationType = utils.Ptr("other")
				m.DestinationValue = utils.Ptr("1.1.1.1")
			}),
			expected: nil,
		},
		{
			description: "nil destination value",
			model: fixtureInputModel(func(m *inputModel) {
				m.DestinationValue = nil
			}),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := buildDestination(tt.model)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("buildDestination() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiAddRoutesToRoutingTableRequest
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
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutesToRoutingTableRequest) {
				*request = (*request).AddRoutesToRoutingTablePayload(fixturePayload(func(payload *iaas.AddRoutesToRoutingTablePayload) {
					(*payload.Items)[0].Labels = utils.ConvertStringMapToInterfaceMap(utils.Ptr(map[string]string{"key": "value"}))
				}))
			}),
		},
		{
			description: "destination is cidrv6 and nexthop is ipv6",
			model: fixtureInputModel(func(model *inputModel) {
				model.DestinationType = utils.Ptr("cidrv6")
				model.DestinationValue = utils.Ptr("2001:db8::/32")
				model.NextHopType = utils.Ptr("ipv6")
				model.NextHopValue = utils.Ptr("2001:db8::1")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutesToRoutingTableRequest) {
				*request = (*request).AddRoutesToRoutingTablePayload(iaas.AddRoutesToRoutingTablePayload{
					Items: &[]iaas.Route{
						{
							Destination: &iaas.RouteDestination{
								DestinationCIDRv6: &iaas.DestinationCIDRv6{
									Type:  utils.Ptr("cidrv6"),
									Value: utils.Ptr("2001:db8::/32"),
								},
							},
							Nexthop: &iaas.RouteNexthop{
								NexthopIPv6: &iaas.NexthopIPv6{
									Type:  utils.Ptr("ipv6"),
									Value: utils.Ptr("2001:db8::1"),
								},
							},
							Labels: utils.ConvertStringMapToInterfaceMap(testLabels),
						},
					},
				})
			}),
		},
		{
			description: "nexthop type is internet (no value)",
			model: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = utils.Ptr("internet")
				model.NextHopValue = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutesToRoutingTableRequest) {
				payload := fixturePayload(func(payload *iaas.AddRoutesToRoutingTablePayload) {
					(*payload.Items)[0].Nexthop = &iaas.RouteNexthop{
						NexthopInternet: &iaas.NexthopInternet{
							Type: utils.Ptr("internet"),
						},
					}
				})
				*request = (*request).AddRoutesToRoutingTablePayload(payload)
			}),
		},
		{
			description: "nexthop type is blackhole (no value)",
			model: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = utils.Ptr("blackhole")
				model.NextHopValue = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutesToRoutingTableRequest) {
				payload := fixturePayload(func(payload *iaas.AddRoutesToRoutingTablePayload) {
					(*payload.Items)[0].Nexthop = &iaas.RouteNexthop{
						NexthopBlackhole: &iaas.NexthopBlackhole{
							Type: utils.Ptr("blackhole"),
						},
					}
				})
				*request = (*request).AddRoutesToRoutingTablePayload(payload)
			}),
		},
		{
			description: "nexthop type is ipv4 with value",
			model: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = utils.Ptr("ipv4")
				model.NextHopValue = utils.Ptr("1.2.3.4")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutesToRoutingTableRequest) {
				payload := fixturePayload(func(payload *iaas.AddRoutesToRoutingTablePayload) {
					(*payload.Items)[0].Nexthop = &iaas.RouteNexthop{
						NexthopIPv4: &iaas.NexthopIPv4{
							Type:  utils.Ptr("ipv4"),
							Value: utils.Ptr("1.2.3.4"),
						},
					}
				})
				*request = (*request).AddRoutesToRoutingTablePayload(payload)
			}),
		},
		{
			description: "nexthop type is ipv6 with value",
			model: fixtureInputModel(func(model *inputModel) {
				model.NextHopType = utils.Ptr("ipv6")
				model.NextHopValue = utils.Ptr("2001:db8::1")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiAddRoutesToRoutingTableRequest) {
				payload := fixturePayload(func(payload *iaas.AddRoutesToRoutingTablePayload) {
					(*payload.Items)[0].Nexthop = &iaas.RouteNexthop{
						NexthopIPv6: &iaas.NexthopIPv6{
							Type:  utils.Ptr("ipv6"),
							Value: utils.Ptr("2001:db8::1"),
						},
					}
				})
				*request = (*request).AddRoutesToRoutingTablePayload(payload)
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
	dummyRoute := iaas.Route{
		Id: utils.Ptr("route-foo"),
		Destination: &iaas.RouteDestination{
			DestinationCIDRv4: &iaas.DestinationCIDRv4{
				Type:  utils.Ptr("cidrv4"),
				Value: utils.Ptr("10.0.0.0/24"),
			},
		},
		Nexthop: &iaas.RouteNexthop{
			NexthopIPv4: &iaas.NexthopIPv4{
				Type:  utils.Ptr("ipv4"),
				Value: utils.Ptr("10.0.0.1"),
			},
		},
		Labels:    utils.ConvertStringMapToInterfaceMap(testLabels),
		CreatedAt: utils.Ptr(time.Now()),
		UpdatedAt: utils.Ptr(time.Now()),
	}

	tests := []struct {
		name         string
		outputFormat string
		items        []iaas.Route
		wantErr      bool
	}{
		{
			name:         "nil items should return error",
			outputFormat: "",
			items:        nil,
			wantErr:      true,
		},
		{
			name:         "empty items list",
			outputFormat: "",
			items:        []iaas.Route{},
			wantErr:      true,
		},
		{
			name:         "table output with one route",
			outputFormat: "",
			items:        []iaas.Route{dummyRoute},
			wantErr:      false,
		},
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			items:        []iaas.Route{dummyRoute},
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			items:        []iaas.Route{dummyRoute},
			wantErr:      false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.outputFormat, tt.items); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
