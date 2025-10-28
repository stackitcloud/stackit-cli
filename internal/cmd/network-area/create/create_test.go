package create

import (
	"context"
	"strconv"
	"strings"
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
	testRegion                    = "eu01"
	testName                      = "example-network-area-name"
	testTransferNetwork           = "100.0.0.0/24"
	testDefaultPrefixLength int64 = 25
	testMaxPrefixLength     int64 = 26
	testMinPrefixLength     int64 = 24
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var (
	testOrgId          = uuid.NewString()
	testAreaId         = uuid.NewString()
	testDnsNameservers = []string{"1.1.1.0", "1.1.2.0"}
	testNetworkRanges  = []string{"192.0.0.0/24", "102.0.0.0/24"}
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,

		nameFlag:           testName,
		organizationIdFlag: testOrgId,
		labelFlag:          "key=value",
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
		Name:           utils.Ptr("example-network-area-name"),
		OrganizationId: testOrgId,
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
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureRequestRegionalArea(mods ...func(request *iaas.ApiCreateNetworkAreaRegionRequest)) iaas.ApiCreateNetworkAreaRegionRequest {
	req := testClient.CreateNetworkAreaRegion(testCtx, testOrgId, testAreaId, testRegion)
	req = req.CreateNetworkAreaRegionPayload(fixtureRegionalAreaPayload())
	for _, mod := range mods {
		mod(&req)
	}
	return req
}

func fixtureRegionalAreaPayload(mods ...func(request *iaas.CreateNetworkAreaRegionPayload)) iaas.CreateNetworkAreaRegionPayload {
	var networkRanges []iaas.NetworkRange
	for _, networkRange := range testNetworkRanges {
		networkRanges = append(networkRanges, iaas.NetworkRange{
			Prefix: utils.Ptr(networkRange),
		})
	}

	payload := iaas.CreateNetworkAreaRegionPayload{
		Ipv4: &iaas.RegionalAreaIPv4{
			DefaultNameservers: utils.Ptr(testDnsNameservers),
			DefaultPrefixLen:   utils.Ptr(testDefaultPrefixLength),
			MaxPrefixLen:       utils.Ptr(testMaxPrefixLength),
			MinPrefixLen:       utils.Ptr(testMinPrefixLength),
			NetworkRanges:      utils.Ptr(networkRanges),
			TransferNetwork:    utils.Ptr(testTransferNetwork),
		},
		Status: nil,
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
			description: "with deprecated flags",
			flagValues: map[string]string{
				nameFlag:           testName,
				organizationIdFlag: testOrgId,

				// Deprecated flags
				dnsNameServersFlag:      strings.Join(testDnsNameservers, ","),
				networkRangesFlag:       strings.Join(testNetworkRanges, ","),
				transferNetworkFlag:     testTransferNetwork,
				defaultPrefixLengthFlag: strconv.FormatInt(testDefaultPrefixLength, 10),
				maxPrefixLengthFlag:     strconv.FormatInt(testMaxPrefixLength, 10),
				minPrefixLengthFlag:     strconv.FormatInt(testMinPrefixLength, 10),
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					Verbosity: globalflags.VerbosityDefault,
				},
				Name:           utils.Ptr(testName),
				OrganizationId: testOrgId,

				// Deprecated fields
				DnsNameServers:      utils.Ptr(testDnsNameservers),
				NetworkRanges:       utils.Ptr(testNetworkRanges),
				TransferNetwork:     utils.Ptr(testTransferNetwork),
				DefaultPrefixLength: utils.Ptr(testDefaultPrefixLength),
				MaxPrefixLength:     utils.Ptr(testMaxPrefixLength),
				MinPrefixLength:     utils.Ptr(testMinPrefixLength),
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
			description: "deprecated network ranges missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkRangesFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NetworkRanges = nil
			}),
		},
		{
			description: "deprecated transfer network missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, transferNetworkFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.TransferNetwork = nil
			}),
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

func TestBuildRequestNetworkAreaRegion(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		areaId          string
		expectedRequest iaas.ApiCreateNetworkAreaRegionRequest
	}{
		{
			description: "base",
			model: fixtureInputModel(func(model *inputModel) {
				// Deprecated fields
				model.DnsNameServers = utils.Ptr(testDnsNameservers)
				model.NetworkRanges = utils.Ptr(testNetworkRanges)
				model.TransferNetwork = utils.Ptr(testTransferNetwork)
				model.DefaultPrefixLength = utils.Ptr(testDefaultPrefixLength)
				model.MaxPrefixLength = utils.Ptr(testMaxPrefixLength)
				model.MinPrefixLength = utils.Ptr(testMinPrefixLength)
			}),
			areaId:          testAreaId,
			expectedRequest: fixtureRequestRegionalArea(),
		},
		{
			description: "base without network ranges",
			model: fixtureInputModel(func(model *inputModel) {
				// Deprecated fields
				model.DnsNameServers = utils.Ptr(testDnsNameservers)
				model.NetworkRanges = utils.Ptr(testNetworkRanges)
				model.TransferNetwork = utils.Ptr(testTransferNetwork)
				model.DefaultPrefixLength = utils.Ptr(testDefaultPrefixLength)
				model.MaxPrefixLength = utils.Ptr(testMaxPrefixLength)
				model.MinPrefixLength = utils.Ptr(testMinPrefixLength)
			}),
			areaId:          testAreaId,
			expectedRequest: fixtureRequestRegionalArea(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequestNetworkAreaRegion(testCtx, tt.model, testAreaId, testClient)

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
		responses    *NetworkAreaResponses
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
			name: "set empty response",
			args: args{
				responses: &NetworkAreaResponses{},
			},
			wantErr: false,
		},
		{
			name: "set empty network area",
			args: args{
				responses: &NetworkAreaResponses{
					NetworkArea: iaas.NetworkArea{},
				},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.orgLabel, tt.args.responses); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfiguredDeprecatedFlags(t *testing.T) {
	type args struct {
		model *inputModel
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no deprecated flags",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{
						Verbosity: globalflags.VerbosityDefault,
					},
					Name:           utils.Ptr(testName),
					OrganizationId: testOrgId,
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      nil,
					NetworkRanges:       nil,
					TransferNetwork:     nil,
					DefaultPrefixLength: nil,
					MaxPrefixLength:     nil,
					MinPrefixLength:     nil,
				},
			},
			want: nil,
		},
		{
			name: "deprecated flags",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{
						Verbosity: globalflags.VerbosityDefault,
					},
					Name:           utils.Ptr(testName),
					OrganizationId: testOrgId,
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      utils.Ptr(testDnsNameservers),
					NetworkRanges:       utils.Ptr(testNetworkRanges),
					TransferNetwork:     utils.Ptr(testTransferNetwork),
					DefaultPrefixLength: utils.Ptr(testDefaultPrefixLength),
					MaxPrefixLength:     utils.Ptr(testMaxPrefixLength),
					MinPrefixLength:     utils.Ptr(testMinPrefixLength),
				},
			},
			want: []string{dnsNameServersFlag, networkRangesFlag, transferNetworkFlag, defaultPrefixLengthFlag, minPrefixLengthFlag, maxPrefixLengthFlag},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getConfiguredDeprecatedFlags(tt.args.model)

			less := func(a, b string) bool {
				return a < b
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(less)); diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestHasDeprecatedFlagsSet(t *testing.T) {
	type args struct {
		model *inputModel
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no deprecated flags",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{
						Verbosity: globalflags.VerbosityDefault,
					},
					Name:           utils.Ptr(testName),
					OrganizationId: testOrgId,
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      nil,
					NetworkRanges:       nil,
					TransferNetwork:     nil,
					DefaultPrefixLength: nil,
					MaxPrefixLength:     nil,
					MinPrefixLength:     nil,
				},
			},
			want: false,
		},
		{
			name: "deprecated flags",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{
						Verbosity: globalflags.VerbosityDefault,
					},
					Name:           utils.Ptr(testName),
					OrganizationId: testOrgId,
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      utils.Ptr(testDnsNameservers),
					NetworkRanges:       utils.Ptr(testNetworkRanges),
					TransferNetwork:     utils.Ptr(testTransferNetwork),
					DefaultPrefixLength: utils.Ptr(testDefaultPrefixLength),
					MaxPrefixLength:     utils.Ptr(testMaxPrefixLength),
					MinPrefixLength:     utils.Ptr(testMinPrefixLength),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasDeprecatedFlagsSet(tt.args.model); got != tt.want {
				t.Errorf("hasDeprecatedFlagsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
