package create

import (
	"context"
	"strconv"
	"strings"
	"testing"

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

const (
	testRegion                    = "eu01"
	testDefaultPrefixLength int64 = 25
	testMaxPrefixLength     int64 = 29
	testMinPrefixLength     int64 = 24
	testTransferNetwork           = "192.168.2.0/24"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var (
	testAreaId             = uuid.NewString()
	testOrgId              = uuid.NewString()
	testDefaultNameservers = []string{"8.8.8.8", "8.8.4.4"}
	testNetworkRanges      = []string{"192.168.0.0/24", "10.0.0.0/24"}
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,

		networkAreaIdFlag:           testAreaId,
		organizationIdFlag:          testOrgId,
		ipv4DefaultNameservers:      strings.Join(testDefaultNameservers, ","),
		ipv4DefaultPrefixLengthFlag: strconv.FormatInt(testDefaultPrefixLength, 10),
		ipv4MaxPrefixLengthFlag:     strconv.FormatInt(testMaxPrefixLength, 10),
		ipv4MinPrefixLengthFlag:     strconv.FormatInt(testMinPrefixLength, 10),
		ipv4NetworkRangesFlag:       strings.Join(testNetworkRanges, ","),
		ipv4TransferNetworkFlag:     testTransferNetwork,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		OrganizationId:          testOrgId,
		NetworkAreaId:           testAreaId,
		IPv4DefaultNameservers:  utils.Ptr(testDefaultNameservers),
		IPv4DefaultPrefixLength: utils.Ptr(testDefaultPrefixLength),
		IPv4MaxPrefixLength:     utils.Ptr(testMaxPrefixLength),
		IPv4MinPrefixLength:     utils.Ptr(testMinPrefixLength),
		IPv4NetworkRanges:       testNetworkRanges,
		IPv4TransferNetwork:     testTransferNetwork,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateNetworkAreaRegionRequest)) iaas.ApiCreateNetworkAreaRegionRequest {
	request := testClient.CreateNetworkAreaRegion(testCtx, testOrgId, testAreaId, testRegion)
	request = request.CreateNetworkAreaRegionPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateNetworkAreaRegionPayload)) iaas.CreateNetworkAreaRegionPayload {
	var networkRange []iaas.NetworkRange
	if len(testNetworkRanges) > 0 {
		networkRange = make([]iaas.NetworkRange, len(testNetworkRanges))
		for i := range testNetworkRanges {
			networkRange[i] = iaas.NetworkRange{
				Prefix: utils.Ptr(testNetworkRanges[i]),
			}
		}
	}

	payload := iaas.CreateNetworkAreaRegionPayload{
		Ipv4: &iaas.RegionalAreaIPv4{
			DefaultNameservers: utils.Ptr(testDefaultNameservers),
			DefaultPrefixLen:   utils.Ptr(testDefaultPrefixLength),
			MaxPrefixLen:       utils.Ptr(testMaxPrefixLength),
			MinPrefixLen:       utils.Ptr(testMinPrefixLength),
			NetworkRanges:      utils.Ptr(networkRange),
			TransferNetwork:    utils.Ptr(testTransferNetwork),
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "area id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "area id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "area id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = "invalid-uuid"
			}),
			isValid: false,
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
			description: "network range missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipv4NetworkRangesFlag)
			}),
			isValid: false,
		},
		{
			description: "multiple network ranges",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4NetworkRangesFlag] = "192.168.2.0/24,10.0.0.0/24"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IPv4NetworkRanges = []string{"192.168.2.0/24", "10.0.0.0/24"}
			}),
			isValid: true,
		},
		{
			description: "network range invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4NetworkRangesFlag] = "invalid-cidr"
			}),
			isValid: false,
		},
		{
			description: "transfer network missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipv4TransferNetworkFlag)
			}),
			isValid: false,
		},
		{
			description: "transfer network invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4TransferNetworkFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "transfer network invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipv4TransferNetworkFlag] = "invalid-cidr"
			}),
			isValid: false,
		},
		{
			description: "region empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.RegionFlag] = ""
			}),
			isValid: false,
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
		expectedRequest iaas.ApiCreateNetworkAreaRegionRequest
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
		outputFormat     string
		region           string
		networkAreaLabel string
		regionalArea     iaas.RegionalArea
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
			name: "set empty regional area",
			args: args{
				regionalArea: iaas.RegionalArea{},
			},
			wantErr: false,
		},
		{
			name: "output json",
			args: args{
				outputFormat: print.JSONOutputFormat,
				regionalArea: iaas.RegionalArea{},
			},
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.region, tt.args.networkAreaLabel, tt.args.regionalArea); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
