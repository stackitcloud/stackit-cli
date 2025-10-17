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

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testOrgId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		nameFlag:                "example-network-area-name",
		organizationIdFlag:      testOrgId,
		dnsNameServersFlag:      "1.1.1.0,1.1.2.0",
		networkRangesFlag:       "192.0.0.0/24,102.0.0.0/24",
		transferNetworkFlag:     "100.0.0.0/24",
		defaultPrefixLengthFlag: "24",
		maxPrefixLengthFlag:     "24",
		minPrefixLengthFlag:     "24",
		labelFlag:               "key=value",
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
		Name:                utils.Ptr("example-network-area-name"),
		OrganizationId:      utils.Ptr(testOrgId),
		DnsNameServers:      utils.Ptr([]string{"1.1.1.0", "1.1.2.0"}),
		NetworkRanges:       utils.Ptr([]string{"192.0.0.0/24", "102.0.0.0/24"}),
		TransferNetwork:     utils.Ptr("100.0.0.0/24"),
		DefaultPrefixLength: utils.Ptr(int64(24)),
		MaxPrefixLength:     utils.Ptr(int64(24)),
		MinPrefixLength:     utils.Ptr(int64(24)),
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateNetworkAreaRequest)) iaas.ApiCreateNetworkAreaRequest {
	request := testClient.CreateNetworkArea(testCtx, testOrgId)
	request = request.CreateNetworkAreaPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateNetworkAreaPayload)) iaas.CreateNetworkAreaPayload {
	payload := iaas.CreateNetworkAreaPayload{
		Name: utils.Ptr("example-network-area-name"),
		Labels: utils.Ptr(map[string]interface{}{
			"key": "value",
		}),
		AddressFamily: &iaas.CreateAreaAddressFamily{
			Ipv4: &iaas.CreateAreaIPv4{
				DefaultNameservers: utils.Ptr([]string{"1.1.1.0", "1.1.2.0"}),
				NetworkRanges: &[]iaas.NetworkRange{
					{
						Prefix: utils.Ptr("192.0.0.0/24"),
					},
					{
						Prefix: utils.Ptr("102.0.0.0/24"),
					},
				},
				TransferNetwork:  utils.Ptr("100.0.0.0/24"),
				DefaultPrefixLen: utils.Ptr(int64(24)),
				MaxPrefixLen:     utils.Ptr(int64(24)),
				MinPrefixLen:     utils.Ptr(int64(24)),
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
			description: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, dnsNameServersFlag)
				delete(flagValues, defaultPrefixLengthFlag)
				delete(flagValues, maxPrefixLengthFlag)
				delete(flagValues, minPrefixLengthFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DnsNameServers = nil
				model.DefaultPrefixLength = nil
				model.MaxPrefixLength = nil
				model.MinPrefixLength = nil
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
			description: "network ranges missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkRangesFlag)
			}),
			isValid: false,
		},
		{
			description: "transfer network missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, transferNetworkFlag)
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
			description: "org id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
			}),
			isValid: false,
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
		expectedRequest iaas.ApiCreateNetworkAreaRequest
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

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat string
		orgLabel     string
		networkArea  *iaas.NetworkArea
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
			name: "set empty network area",
			args: args{
				networkArea: &iaas.NetworkArea{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.orgLabel, tt.args.networkArea); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
