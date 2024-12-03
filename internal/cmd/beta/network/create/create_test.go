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

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testProjectId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		nameFlag:               "example-network-name",
		ipv4DnsNameServersFlag: "1.1.1.0,1.1.2.0",
		ipv4PrefixLengthFlag:   "24",
		ipv4PrefixFlag:         "10.1.2.0/24",
		ipv4GatewayFlag:        "10.1.2.3",
		ipv6DnsNameServersFlag: "2001:4860:4860::8888,2001:4860:4860::8844",
		ipv6PrefixLengthFlag:   "24",
		ipv6PrefixFlag:         "2001:4860:4860::8888",
		ipv6GatewayFlag:        "2001:4860:4860::8888",
		routedFlag:             "true",
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
			Verbosity: globalflags.VerbosityDefault,
		},
		Name:               utils.Ptr("example-network-name"),
		IPv4DnsNameServers: utils.Ptr([]string{"1.1.1.0", "1.1.2.0"}),
		IPv4PrefixLength:   utils.Ptr(int64(24)),
		IPv4Prefix:         utils.Ptr("10.1.2.0/24"),
		IPv4Gateway:        utils.Ptr("10.1.2.3"),
		IPv6DnsNameServers: utils.Ptr([]string{"2001:4860:4860::8888", "2001:4860:4860::8844"}),
		IPv6PrefixLength:   utils.Ptr(int64(24)),
		IPv6Prefix:         utils.Ptr("2001:4860:4860::8888"),
		IPv6Gateway:        utils.Ptr("2001:4860:4860::8888"),
		Routed:             true,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateNetworkRequest)) iaas.ApiCreateNetworkRequest {
	request := testClient.CreateNetwork(testCtx, testProjectId)
	request = request.CreateNetworkPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequiredRequest(mods ...func(request *iaas.ApiCreateNetworkRequest)) iaas.ApiCreateNetworkRequest {
	request := testClient.CreateNetwork(testCtx, testProjectId)
	request = request.CreateNetworkPayload(iaas.CreateNetworkPayload{
		Name:   utils.Ptr("example-network-name"),
		Routed: utils.Ptr(false),
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
		AddressFamily: &iaas.CreateNetworkAddressFamily{
			Ipv4: &iaas.CreateNetworkIPv4Body{
				Nameservers:  utils.Ptr([]string{"1.1.1.0", "1.1.2.0"}),
				PrefixLength: utils.Ptr(int64(24)),
				Prefix:       utils.Ptr("10.1.2.0/24"),
				Gateway:      iaas.NewNullableString(utils.Ptr("10.1.2.3")),
			},
			Ipv6: &iaas.CreateNetworkIPv6Body{
				Nameservers:  utils.Ptr([]string{"2001:4860:4860::8888", "2001:4860:4860::8844"}),
				PrefixLength: utils.Ptr(int64(24)),
				Prefix:       utils.Ptr("2001:4860:4860::8888"),
				Gateway:      iaas.NewNullableString(utils.Ptr("2001:4860:4860::8888")),
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
			description: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipv4DnsNameServersFlag)
				delete(flagValues, ipv4PrefixLengthFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv4DnsNameServers = nil
				model.IPv4PrefixLength = nil
			}),
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
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "use dns servers, prefix, gateway and prefix length",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4DnsNameServersFlag] = "1.1.1.1"
				flagValues[ipv4PrefixLengthFlag] = "25"
				flagValues[ipv4PrefixFlag] = "10.1.2.0/24"
				flagValues[ipv4GatewayFlag] = "10.1.2.3"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv4DnsNameServers = utils.Ptr([]string{"1.1.1.1"})
				model.IPv4PrefixLength = utils.Ptr(int64(25))
				model.IPv4Prefix = utils.Ptr("10.1.2.0/24")
				model.IPv4Gateway = utils.Ptr("10.1.2.3")
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
			description: "use ipv6 dns servers, prefix, gateway and prefix length",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv6DnsNameServersFlag] = "2001:4860:4860::8888"
				flagValues[ipv6PrefixLengthFlag] = "25"
				flagValues[ipv6PrefixFlag] = "2001:4860:4860::8888"
				flagValues[ipv6GatewayFlag] = "2001:4860:4860::8888"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv6DnsNameServers = utils.Ptr([]string{"2001:4860:4860::8888"})
				model.IPv6PrefixLength = utils.Ptr(int64(25))
				model.IPv6Prefix = utils.Ptr("2001:4860:4860::8888")
				model.IPv6Gateway = utils.Ptr("2001:4860:4860::8888")
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
			description: "use routed true",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[routedFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Routed = true
			}),
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
				},
				Name: utils.Ptr("example-network-name"),
			},
			expectedRequest: fixtureRequiredRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
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
