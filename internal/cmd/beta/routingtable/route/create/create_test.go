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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaasalpha"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaasalpha.APIClient{}

var testRegion = "eu01"
var testOrgId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testRoutingTableId = uuid.NewString()

var testDestinationTypeFlag = "cidrv4"
var testDestinationValueFlag = "1.1.1.0/24"
var testNextHopTypeFlag = "ipv4"
var testNextHopValueFlag = "1.1.1.1"
var testLabelSelectorFlag = "key1=value1,key2=value2"
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
		OrganizationId:   utils.Ptr(testOrgId),
		NetworkAreaId:    utils.Ptr(testNetworkAreaId),
		RoutingTableId:   utils.Ptr(testRoutingTableId),
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

func fixtureRequest(mods ...func(request *iaasalpha.ApiAddRoutesToRoutingTableRequest)) iaasalpha.ApiAddRoutesToRoutingTableRequest {
	request := testClient.AddRoutesToRoutingTable(testCtx, testOrgId, testNetworkAreaId, testRegion, testRoutingTableId)
	request = request.AddRoutesToRoutingTablePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaasalpha.AddRoutesToRoutingTablePayload)) iaasalpha.AddRoutesToRoutingTablePayload {
	payload := iaasalpha.AddRoutesToRoutingTablePayload{
		Items: &[]iaasalpha.Route{
			{
				Destination: &iaasalpha.RouteDestination{
					DestinationCIDRv4: &iaasalpha.DestinationCIDRv4{
						Type:  utils.Ptr(testDestinationTypeFlag),
						Value: utils.Ptr(testDestinationValueFlag),
					},
				},
				Nexthop: &iaasalpha.RouteNexthop{
					NexthopIPv4: &iaasalpha.NexthopIPv4{
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
			description: "routing-table-id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, routingTableIdFlag)
			}),
			isValid: false,
		},
		{
			description: "destination-value missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, destinationValueFlag)
			}),
			isValid: false,
		},
		{
			description: "destination-type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, destinationTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "nexthop-type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nextHopTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "nexthop-value missing",
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
			description: "invalid destination type enum",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationTypeFlag] = "ipv4"
			}),
			isValid: false,
		},
		{
			description: "destination value not ipv4 cidr",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationValueFlag] = "0.0.0.0"
			}),
			isValid: false,
		},
		{
			description: "destination value not ipv6 cidr",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[destinationTypeFlag] = "cidrv6"
				flagValues[destinationValueFlag] = "2001:db8::"
			}),
			isValid: false,
		},
		{
			description: "destination value is ipv6 cidr",
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
			description: "nexthop-type is internet and nexthop-value is provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "internet"
				flagValues[nextHopValueFlag] = "1.1.1.1" // should not be allowed
			}),
			isValid: false,
		},
		{
			description: "nexthop-type is blackhole and nexthop-value is provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "blackhole"
				flagValues[nextHopValueFlag] = "1.1.1.1"
			}),
			isValid: false,
		},
		{
			description: "nexthop-type is internet and nexthop-value is not provided",
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
			description: "nexthop-type is blackhole and nexthop-value is not provided",
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
			description: "nexthop-type is ipv4 and nexthop-value is missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "ipv4"
				delete(flagValues, nextHopValueFlag)
			}),
			isValid: false,
		},
		{
			description: "nexthop-type is ipv6 and nexthop-value is missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "ipv6"
				delete(flagValues, nextHopValueFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid nexthop-type provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nextHopTypeFlag] = "invalid-type"
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

func TestBuildNextHop(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expected    *iaasalpha.RouteNexthop
	}{
		{
			description: "IPv4 next hop",
			model: fixtureInputModel(func(m *inputModel) {
				m.NextHopType = utils.Ptr("ipv4")
				m.NextHopValue = utils.Ptr("1.1.1.1")
			}),
			expected: &iaasalpha.RouteNexthop{
				NexthopIPv4: &iaasalpha.NexthopIPv4{
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
			expected: &iaasalpha.RouteNexthop{
				NexthopIPv6: &iaasalpha.NexthopIPv6{
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
			expected: &iaasalpha.RouteNexthop{
				NexthopInternet: &iaasalpha.NexthopInternet{
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
			expected: &iaasalpha.RouteNexthop{
				NexthopBlackhole: &iaasalpha.NexthopBlackhole{
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
		expected    *iaasalpha.RouteDestination
	}{
		{
			description: "CIDRv4 destination",
			model: fixtureInputModel(func(m *inputModel) {
				m.DestinationType = utils.Ptr("cidrv4")
				m.DestinationValue = utils.Ptr("192.168.1.0/24")
			}),
			expected: &iaasalpha.RouteDestination{
				DestinationCIDRv4: &iaasalpha.DestinationCIDRv4{
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
			expected: &iaasalpha.RouteDestination{
				DestinationCIDRv6: &iaasalpha.DestinationCIDRv6{
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
		expectedRequest iaasalpha.ApiAddRoutesToRoutingTableRequest
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
			expectedRequest: fixtureRequest(func(request *iaasalpha.ApiAddRoutesToRoutingTableRequest) {
				*request = (*request).AddRoutesToRoutingTablePayload(fixturePayload(func(payload *iaasalpha.AddRoutesToRoutingTablePayload) {
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
			expectedRequest: fixtureRequest(func(request *iaasalpha.ApiAddRoutesToRoutingTableRequest) {
				*request = (*request).AddRoutesToRoutingTablePayload(iaasalpha.AddRoutesToRoutingTablePayload{
					Items: &[]iaasalpha.Route{
						{
							Destination: &iaasalpha.RouteDestination{
								DestinationCIDRv6: &iaasalpha.DestinationCIDRv6{
									Type:  utils.Ptr("cidrv6"),
									Value: utils.Ptr("2001:db8::/32"),
								},
							},
							Nexthop: &iaasalpha.RouteNexthop{
								NexthopIPv6: &iaasalpha.NexthopIPv6{
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
			expectedRequest: fixtureRequest(func(request *iaasalpha.ApiAddRoutesToRoutingTableRequest) {
				payload := fixturePayload(func(payload *iaasalpha.AddRoutesToRoutingTablePayload) {
					(*payload.Items)[0].Nexthop = &iaasalpha.RouteNexthop{
						NexthopInternet: &iaasalpha.NexthopInternet{
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
			expectedRequest: fixtureRequest(func(request *iaasalpha.ApiAddRoutesToRoutingTableRequest) {
				payload := fixturePayload(func(payload *iaasalpha.AddRoutesToRoutingTablePayload) {
					(*payload.Items)[0].Nexthop = &iaasalpha.RouteNexthop{
						NexthopBlackhole: &iaasalpha.NexthopBlackhole{
							Type: utils.Ptr("blackhole"),
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
	dummyRoute := iaasalpha.Route{
		Id: utils.Ptr("route-foo"),
		Destination: &iaasalpha.RouteDestination{
			DestinationCIDRv4: &iaasalpha.DestinationCIDRv4{
				Type:  utils.Ptr("cidrv4"),
				Value: utils.Ptr("10.0.0.0/24"),
			},
		},
		Nexthop: &iaasalpha.RouteNexthop{
			NexthopIPv4: &iaasalpha.NexthopIPv4{
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
		items        []iaasalpha.Route
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
			items:        []iaasalpha.Route{},
			wantErr:      true,
		},
		{
			name:         "table output with one route",
			outputFormat: "",
			items:        []iaasalpha.Route{dummyRoute},
			wantErr:      false,
		},
		{
			name:         "json output with one route",
			outputFormat: print.JSONOutputFormat,
			items:        []iaasalpha.Route{dummyRoute},
			wantErr:      false,
		},
		{
			name:         "yaml output with one route",
			outputFormat: print.YAMLOutputFormat,
			items:        []iaasalpha.Route{dummyRoute},
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
