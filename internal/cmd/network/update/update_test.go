package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
var testNetworkId = uuid.NewString()

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testNetworkId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		nameFlag:               "example-network-name",
		projectIdFlag:          testProjectId,
		ipv4DnsNameServersFlag: "1.1.1.0,1.1.2.0",
		ipv4GatewayFlag:        "10.1.2.3",
		ipv6DnsNameServersFlag: "2001:4860:4860::8888,2001:4860:4860::8844",
		ipv6GatewayFlag:        "2001:4860:4860::8888",
		labelFlag:              "key=value",
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
		NetworkId:          testNetworkId,
		IPv4DnsNameServers: utils.Ptr([]string{"1.1.1.0", "1.1.2.0"}),
		IPv4Gateway:        utils.Ptr("10.1.2.3"),
		IPv6DnsNameServers: utils.Ptr([]string{"2001:4860:4860::8888", "2001:4860:4860::8844"}),
		IPv6Gateway:        utils.Ptr("2001:4860:4860::8888"),
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiPartialUpdateNetworkRequest)) iaas.ApiPartialUpdateNetworkRequest {
	request := testClient.PartialUpdateNetwork(testCtx, testProjectId, testNetworkId)
	request = request.PartialUpdateNetworkPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.PartialUpdateNetworkPayload)) iaas.PartialUpdateNetworkPayload {
	payload := iaas.PartialUpdateNetworkPayload{
		Name: utils.Ptr("example-network-name"),
		Labels: utils.Ptr(map[string]interface{}{
			"key": "value",
		}),
		AddressFamily: &iaas.UpdateNetworkAddressFamily{
			Ipv4: &iaas.UpdateNetworkIPv4Body{
				Nameservers: utils.Ptr([]string{"1.1.1.0", "1.1.2.0"}),
				Gateway:     iaas.NewNullableString(utils.Ptr("10.1.2.3")),
			},
			Ipv6: &iaas.UpdateNetworkIPv6Body{
				Nameservers: utils.Ptr([]string{"2001:4860:4860::8888", "2001:4860:4860::8844"}),
				Gateway:     iaas.NewNullableString(utils.Ptr("2001:4860:4860::8888")),
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "required only",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipv4DnsNameServersFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv4DnsNameServers = nil
			}),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "network id invalid 1",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "network id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "use dns servers and gateway",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4DnsNameServersFlag] = "1.1.1.1"
				flagValues[ipv4GatewayFlag] = "10.1.2.3"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv4DnsNameServers = utils.Ptr([]string{"1.1.1.1"})
				model.IPv4Gateway = utils.Ptr("10.1.2.3")
			}),
		},
		{
			description: "use ipv4 gateway nil",
			argValues:   fixtureArgValues(),
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
			description: "use ipv6 dns servers and gateway",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv6DnsNameServersFlag] = "2001:4860:4860::8888"
				flagValues[ipv6GatewayFlag] = "2001:4860:4860::8888"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv6DnsNameServers = utils.Ptr([]string{"2001:4860:4860::8888"})
				model.IPv6Gateway = utils.Ptr("2001:4860:4860::8888")
			}),
		},
		{
			description: "use ipv6 gateway nil",
			argValues:   fixtureArgValues(),
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
			description: "labels missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
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

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
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
		expectedRequest iaas.ApiPartialUpdateNetworkRequest
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
				cmp.AllowUnexported(iaas.NullableString{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
