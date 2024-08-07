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
		ipv6DnsNameServersFlag: "2001:4860:4860::8888,2001:4860:4860::8844",
		ipv6PrefixLengthFlag:   "24",
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
		IPv6DnsNameServers: utils.Ptr([]string{"2001:4860:4860::8888", "2001:4860:4860::8844"}),
		IPv6PrefixLength:   utils.Ptr(int64(24)),
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

func fixturePayload(mods ...func(payload *iaas.CreateNetworkPayload)) iaas.CreateNetworkPayload {
	payload := iaas.CreateNetworkPayload{
		Name: utils.Ptr("example-network-name"),
		AddressFamily: &iaas.CreateNetworkAddressFamily{
			Ipv4: &iaas.CreateNetworkIPv4{
				Nameservers:  utils.Ptr([]string{"1.1.1.0", "1.1.2.0"}),
				PrefixLength: utils.Ptr(int64(24)),
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
			description: "use dns servers and prefix",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4DnsNameServersFlag] = "1.1.1.1"
				flagValues[ipv4PrefixLengthFlag] = "25"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv4DnsNameServers = utils.Ptr([]string{"1.1.1.1"})
				model.IPv4PrefixLength = utils.Ptr(int64(25))
			}),
		},
		{
			description: "use ipv6 dns servers and prefix",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv6DnsNameServersFlag] = "2001:4860:4860::8888"
				flagValues[ipv6PrefixLengthFlag] = "25"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv6DnsNameServers = utils.Ptr([]string{"2001:4860:4860::8888"})
				model.IPv6PrefixLength = utils.Ptr(int64(25))
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
