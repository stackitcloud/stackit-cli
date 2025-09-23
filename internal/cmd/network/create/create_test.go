package create

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

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
	testRegion = "eu01"

	testNetworkName            = "example-network-name"
	testIPv4PrefixLength int64 = 24
	testIPv4Prefix             = "10.1.2.0/24"
	testIPv4Gateway            = "10.1.2.3"
	testIPv6PrefixLength int64 = 24
	testIPv6Prefix             = "2001:4860:4860::/64"
	testIPv6Gateway            = "2001:db8:0:8d3:0:8a2e:70:1"
	testNonRouted              = false
)

var (
	testIPv4NameServers = []string{"1.1.1.0", "1.1.2.0"}
	testIPv6NameServers = []string{"2001:4860:4860::8888", "2001:4860:4860::8844"}
)

type testCtxKey struct{}

var (
	testCtx            = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient         = &iaas.APIClient{}
	testProjectId      = uuid.NewString()
	testRoutingTableId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		nameFlag:           testNetworkName,
		nonRoutedFlag:      strconv.FormatBool(testNonRouted),
		labelFlag:          "key=value",
		routingTableIdFlag: testRoutingTableId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagValuesWithPrefix(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := fixtureFlagValues(func(flagValues map[string]string) {
		flagValues[ipv4DnsNameServersFlag] = strings.Join(testIPv4NameServers, ",")
		flagValues[ipv4PrefixFlag] = testIPv4Prefix
		flagValues[ipv4GatewayFlag] = testIPv4Gateway

		flagValues[ipv6DnsNameServersFlag] = strings.Join(testIPv6NameServers, ",")
		flagValues[ipv6PrefixFlag] = testIPv6Prefix
		flagValues[ipv6GatewayFlag] = testIPv6Gateway
	})
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagValuesWithPrefixLength(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := fixtureFlagValues(func(flagValues map[string]string) {
		flagValues[ipv4PrefixLengthFlag] = strconv.FormatInt(testIPv4PrefixLength, 10)
		flagValues[ipv4DnsNameServersFlag] = strings.Join(testIPv4NameServers, ",")

		flagValues[ipv6PrefixLengthFlag] = strconv.FormatInt(testIPv6PrefixLength, 10)
		flagValues[ipv6DnsNameServersFlag] = strings.Join(testIPv6NameServers, ",")
	})
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		Name:      utils.Ptr(testNetworkName),
		NonRouted: testNonRouted,
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
		RoutingTableID: utils.Ptr(testRoutingTableId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureInputModelWithPrefix(mods ...func(model *inputModel)) *inputModel {
	model := fixtureInputModel()

	model.IPv4DnsNameServers = utils.Ptr(testIPv4NameServers)
	model.IPv4Prefix = utils.Ptr(testIPv4Prefix)
	model.IPv4Gateway = utils.Ptr(testIPv4Gateway)

	model.IPv6DnsNameServers = utils.Ptr(testIPv6NameServers)
	model.IPv6Prefix = utils.Ptr(testIPv6Prefix)
	model.IPv6Gateway = utils.Ptr(testIPv6Gateway)

	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureInputModelWithPrefixLength(mods ...func(model *inputModel)) *inputModel {
	model := fixtureInputModel()

	model.IPv4DnsNameServers = utils.Ptr(testIPv4NameServers)
	model.IPv4PrefixLength = utils.Ptr(testIPv4PrefixLength)

	model.IPv6DnsNameServers = utils.Ptr(testIPv6NameServers)
	model.IPv6PrefixLength = utils.Ptr(testIPv6PrefixLength)

	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateNetworkRequest)) iaas.ApiCreateNetworkRequest {
	request := testClient.CreateNetwork(testCtx, testProjectId, testRegion)
	request = request.CreateNetworkPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequiredRequest(mods ...func(request *iaas.ApiCreateNetworkRequest)) iaas.ApiCreateNetworkRequest {
	request := testClient.CreateNetwork(testCtx, testProjectId, testRegion)
	request = request.CreateNetworkPayload(iaas.CreateNetworkPayload{
		Name:   utils.Ptr(testNetworkName),
		Routed: utils.Ptr(true),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateNetworkPayload)) iaas.CreateNetworkPayload {
	payload := iaas.CreateNetworkPayload{
		Name:   utils.Ptr("example-network-name"),
		Routed: utils.Ptr(true),
		Labels: utils.Ptr(map[string]interface{}{
			"key": "value",
		}),
		RoutingTableId: utils.Ptr(testRoutingTableId),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixturePayloadWithPrefix(mods ...func(payload *iaas.CreateNetworkPayload)) iaas.CreateNetworkPayload {
	payload := fixturePayload()
	payload.Ipv4 = &iaas.CreateNetworkIPv4{
		CreateNetworkIPv4WithPrefix: &iaas.CreateNetworkIPv4WithPrefix{
			Gateway:     iaas.NewNullableString(utils.Ptr(testIPv4Gateway)),
			Nameservers: utils.Ptr(testIPv4NameServers),
			Prefix:      utils.Ptr(testIPv4Prefix),
		},
	}
	payload.Ipv6 = &iaas.CreateNetworkIPv6{
		CreateNetworkIPv6WithPrefix: &iaas.CreateNetworkIPv6WithPrefix{
			Nameservers: utils.Ptr(testIPv6NameServers),
			Prefix:      utils.Ptr(testIPv6Prefix),
			Gateway:     iaas.NewNullableString(utils.Ptr(testIPv6Gateway)),
		},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixturePayloadWithPrefixLength(mods ...func(payload *iaas.CreateNetworkPayload)) iaas.CreateNetworkPayload {
	payload := fixturePayload()
	payload.Ipv4 = &iaas.CreateNetworkIPv4{
		CreateNetworkIPv4WithPrefixLength: &iaas.CreateNetworkIPv4WithPrefixLength{
			PrefixLength: utils.Ptr(testIPv4PrefixLength),
			Nameservers:  utils.Ptr(testIPv4NameServers),
		},
	}
	payload.Ipv6 = &iaas.CreateNetworkIPv6{
		CreateNetworkIPv6WithPrefixLength: &iaas.CreateNetworkIPv6WithPrefixLength{
			PrefixLength: utils.Ptr(testIPv6PrefixLength),
			Nameservers:  utils.Ptr(testIPv6NameServers),
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
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "required only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,

				nameFlag: testNetworkName,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				Name: utils.Ptr(testNetworkName),
			},
		},
		{
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description:   "use with prefix",
			flagValues:    fixtureFlagValuesWithPrefix(),
			isValid:       true,
			expectedModel: fixtureInputModelWithPrefix(),
		},
		{
			description: "use with prefix only ipv4",
			flagValues: fixtureFlagValuesWithPrefix(func(flagValues map[string]string) {
				delete(flagValues, ipv6GatewayFlag)
				delete(flagValues, ipv6PrefixFlag)
				delete(flagValues, ipv6PrefixLengthFlag)
				delete(flagValues, ipv6DnsNameServersFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModelWithPrefix(func(model *inputModel) {
				model.IPv6PrefixLength = nil
				model.IPv6Prefix = nil
				model.IPv6DnsNameServers = nil
				model.IPv6Gateway = nil
			}),
		},
		{
			description: "use with prefix only ipv6",
			flagValues: fixtureFlagValuesWithPrefix(func(flagValues map[string]string) {
				delete(flagValues, ipv4GatewayFlag)
				delete(flagValues, ipv4PrefixFlag)
				delete(flagValues, ipv4PrefixLengthFlag)
				delete(flagValues, ipv4DnsNameServersFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModelWithPrefix(func(model *inputModel) {
				model.IPv4PrefixLength = nil
				model.IPv4Prefix = nil
				model.IPv4DnsNameServers = nil
				model.IPv4Gateway = nil
			}),
		},
		{
			description:   "use with prefixLength",
			flagValues:    fixtureFlagValuesWithPrefixLength(),
			isValid:       true,
			expectedModel: fixtureInputModelWithPrefixLength(),
		},
		{
			description: "use with prefixLength only ipv4",
			flagValues: fixtureFlagValuesWithPrefixLength(func(flagValues map[string]string) {
				delete(flagValues, ipv6GatewayFlag)
				delete(flagValues, ipv6PrefixFlag)
				delete(flagValues, ipv6PrefixLengthFlag)
				delete(flagValues, ipv6DnsNameServersFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModelWithPrefixLength(func(model *inputModel) {
				model.IPv6PrefixLength = nil
				model.IPv6Prefix = nil
				model.IPv6DnsNameServers = nil
				model.IPv6Gateway = nil
			}),
		},
		{
			description: "use with prefixLength only ipv6",
			flagValues: fixtureFlagValuesWithPrefixLength(func(flagValues map[string]string) {
				delete(flagValues, ipv4GatewayFlag)
				delete(flagValues, ipv4PrefixFlag)
				delete(flagValues, ipv4PrefixLengthFlag)
				delete(flagValues, ipv4DnsNameServersFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModelWithPrefixLength(func(model *inputModel) {
				model.IPv4PrefixLength = nil
				model.IPv4Prefix = nil
				model.IPv4DnsNameServers = nil
				model.IPv4Gateway = nil
			}),
		},
		{
			description: "use ipv4 gateway nil",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[noIpv4GatewayFlag] = "true"
				delete(flagValues, ipv4GatewayFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NoIPv4Gateway = true
				model.IPv4Gateway = nil
			}),
		},
		{
			description: "use ipv6 gateway nil",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[noIpv6GatewayFlag] = "true"
				delete(flagValues, ipv6GatewayFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NoIPv6Gateway = true
				model.IPv6Gateway = nil
			}),
		},
		{
			description: "ipv4 prefix length and prefix conflict",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4PrefixFlag] = testIPv4Prefix
				flagValues[ipv4PrefixLengthFlag] = strconv.FormatInt(testIPv4PrefixLength, 10)
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "ipv6 prefix length and prefix conflict",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv6PrefixFlag] = testIPv6Prefix
				flagValues[ipv6PrefixLengthFlag] = strconv.FormatInt(testIPv6PrefixLength, 10)
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "ipv4 nameserver with missing prefix or prefix length",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4DnsNameServersFlag] = strings.Join(testIPv4NameServers, ",")
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "ipv6 nameserver with missing prefix or prefix length",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv6DnsNameServersFlag] = strings.Join(testIPv6NameServers, ",")
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "ipv4 gateway and no-gateway flag conflict",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4GatewayFlag] = testIPv4Gateway
				flagValues[noIpv4GatewayFlag] = "true"
			}),
			isValid: false,
		},
		{
			description: "ipv6 gateway and no-gateway flag conflict",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv6GatewayFlag] = testIPv4Gateway
				flagValues[noIpv6GatewayFlag] = "true"
			}),
			isValid: false,
		},
		{
			description: "ipv4 gateway and prefixLength flag conflict",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4GatewayFlag] = testIPv4Gateway
				flagValues[ipv4PrefixLengthFlag] = strconv.FormatInt(testIPv4PrefixLength, 10)
			}),
			isValid: false,
		},
		{
			description: "ipv6 gateway and prefixLength flag conflict",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv6GatewayFlag] = testIPv6Gateway
				flagValues[ipv6PrefixLengthFlag] = strconv.FormatInt(testIPv6PrefixLength, 10)
			}),
			isValid: false,
		},
		{
			description: "non-routed network",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nonRoutedFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NonRouted = true
			}),
		},
		{
			description: "labels missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
			isValid: true,
		},
		{
			description: "routing-table id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[routingTableIdFlag] = "invalid-uuid"
			}),
			expectedModel: nil,
			isValid:       false,
		},
		{
			description: "routing-table id not set",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, routingTableIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.RoutingTableID = nil
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
	var tests = []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiCreateNetworkRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "only name in payload",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				Name: utils.Ptr(testNetworkName),
			},
			expectedRequest: fixtureRequiredRequest(),
		},
		{
			description: "use prefix length",
			model:       fixtureInputModelWithPrefixLength(),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateNetworkRequest) {
				*request = (*request).CreateNetworkPayload(fixturePayloadWithPrefixLength())
			}),
		},
		{
			description: "use prefix",
			model:       fixtureInputModelWithPrefix(),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateNetworkRequest) {
				*request = (*request).CreateNetworkPayload(fixturePayloadWithPrefix())
			}),
		},
		{
			description: "non-routed network",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				Name:      utils.Ptr(testNetworkName),
				NonRouted: true,
			},
			expectedRequest: testClient.CreateNetwork(testCtx, testProjectId, testRegion).CreateNetworkPayload(iaas.CreateNetworkPayload{
				Name:   utils.Ptr(testNetworkName),
				Routed: utils.Ptr(false),
			}),
		},
		{
			description: "network with routing-table id attached",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				Name:           utils.Ptr(testNetworkName),
				RoutingTableID: utils.Ptr(testRoutingTableId),
			},
			expectedRequest: testClient.CreateNetwork(testCtx, testProjectId, testRegion).CreateNetworkPayload(iaas.CreateNetworkPayload{
				Name:           utils.Ptr(testNetworkName),
				RoutingTableId: utils.Ptr(testRoutingTableId),
				Routed:         utils.Ptr(true),
			}),
		},
		{
			description: "use ipv4 dns servers and prefix length",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				IPv4DnsNameServers: utils.Ptr([]string{"1.1.1.1"}),
				IPv4PrefixLength:   utils.Ptr(int64(25)),
			},
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateNetworkRequest) {
				*request = (*request).CreateNetworkPayload(iaas.CreateNetworkPayload{
					Ipv4: &iaas.CreateNetworkIPv4{
						CreateNetworkIPv4WithPrefixLength: &iaas.CreateNetworkIPv4WithPrefixLength{
							Nameservers:  utils.Ptr([]string{"1.1.1.1"}),
							PrefixLength: utils.Ptr(int64(25)),
						},
					},
					Routed: utils.Ptr(true),
				})
			}),
		},
		{
			description: "use prefix with no gateway",
			model: fixtureInputModelWithPrefix(func(model *inputModel) {
				model.NoIPv4Gateway = true
				model.NoIPv6Gateway = true
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateNetworkRequest) {
				*request = (*request).CreateNetworkPayload(
					fixturePayloadWithPrefix(func(payload *iaas.CreateNetworkPayload) {
						payload.Ipv4.CreateNetworkIPv4WithPrefix.Gateway = iaas.NewNullableString(nil)
						payload.Ipv6.CreateNetworkIPv6WithPrefix.Gateway = iaas.NewNullableString(nil)
					}),
				)
			}),
		},
		{
			description: "use ipv6 dns servers, prefix and gateway",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				IPv6DnsNameServers: utils.Ptr([]string{"2001:4860:4860::8888"}),
				IPv6Prefix:         utils.Ptr("2001:4860:4860::8888"),
				IPv6Gateway:        utils.Ptr("2001:4860:4860::8888"),
			},
			expectedRequest: testClient.CreateNetwork(testCtx, testProjectId, testRegion).CreateNetworkPayload(iaas.CreateNetworkPayload{
				Ipv6: &iaas.CreateNetworkIPv6{
					CreateNetworkIPv6WithPrefix: &iaas.CreateNetworkIPv6WithPrefix{
						Nameservers: utils.Ptr([]string{"2001:4860:4860::8888"}),
						Prefix:      utils.Ptr("2001:4860:4860::8888"),
						Gateway:     iaas.NewNullableString(utils.Ptr("2001:4860:4860::8888")),
					},
				},
				Routed: utils.Ptr(true),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmp.AllowUnexported(iaas.NullableString{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		async        bool
		projectLabel string
		network      *iaas.Network
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "set empty network",
			args: args{
				network: &iaas.Network{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.async, tt.args.projectLabel, tt.args.network); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
