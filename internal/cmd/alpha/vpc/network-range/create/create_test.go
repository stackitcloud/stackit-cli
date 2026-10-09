package create

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	testRegion                    = "eu01"
	testPrefix                    = "10.0.10.0/24"
	testDescription               = "test-description"
	testLabelString               = "env=qa,endpoint=alpha"
	testDefaultPrefixLength int64 = 25
	testMaxPrefixLength     int64 = 29
	testMinPrefixLength     int64 = 24
	testIpVersion                 = iaas.NETWORKRANGEIPV4REQUESTIPVERSION_IPV4
)

type testCtxKey struct{}

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &iaas.APIClient{DefaultAPI: &iaas.DefaultAPIService{}}

	testProjectId   = uuid.NewString()
	testVpcId       = uuid.NewString()
	testNameservers = []string{"1.1.1.1", "8.8.8.8"}
	testLabel       = map[string]any{
		"env":      "qa",
		"endpoint": "alpha",
	}
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag:    testRegion,
		globalflags.ProjectIdFlag: testProjectId,

		vpcIdFlag:               testVpcId,
		prefixFlag:              testPrefix,
		descriptionFlag:         testDescription,
		defaultPrefixLengthFlag: strconv.Itoa(int(testDefaultPrefixLength)),
		maxPrefixLengthFlag:     strconv.Itoa(int(testMaxPrefixLength)),
		minPrefixLengthFlag:     strconv.Itoa(int(testMinPrefixLength)),
		nameserversFlag:         strings.Join(testNameservers, ","),
		labelsFlag:              testLabelString,
		ipVersionFlag.Name():    string(testIpVersion),
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
			ProjectId: testProjectId,
		},
		VpcId:               testVpcId,
		Prefix:              testPrefix,
		IpVersion:           testIpVersion,
		Description:         new(testDescription),
		DefaultPrefixLength: new(testDefaultPrefixLength),
		MaxPrefixLength:     new(testMaxPrefixLength),
		MinPrefixLength:     new(testMinPrefixLength),
		Nameservers:         testNameservers,
		Labels:              testLabel,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateVPCNetworkRangeRequest)) iaas.ApiCreateVPCNetworkRangeRequest {
	request := testClient.DefaultAPI.CreateVPCNetworkRange(testCtx, testProjectId, testVpcId, testRegion)
	request = request.CreateVPCNetworkRangePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateVPCNetworkRangePayload)) iaas.CreateVPCNetworkRangePayload {
	payload := iaas.CreateVPCNetworkRangePayload{
		NetworkRangeIPv4Request: &iaas.NetworkRangeIPv4Request{
			DefaultPrefixLen: new(testDefaultPrefixLength),
			Description:      new(testDescription),
			IpVersion:        testIpVersion,
			Labels:           testLabel,
			MaxPrefixLen:     new(testMaxPrefixLength),
			MinPrefixLen:     new(testMinPrefixLength),
			Nameservers:      testNameservers,
			Prefix:           testPrefix,
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
			description: "prefix missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, prefixFlag)
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
			description: "vpc id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, vpcIdFlag)
			}),
			isValid: false,
		},
		{
			description: "vpc id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[vpcIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "vpc id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[vpcIdFlag] = "invalid-uuid"
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
		expectedRequest iaas.ApiCreateVPCNetworkRangeRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				t.Fatalf("test failed: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, iaas.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	testNetworkRangeId := uuid.NewString()

	type args struct {
		outputFormat string
		vpcLabel     string
		networkRange *iaas.VPCNetworkRange
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
			name: "empty network range",
			args: args{
				networkRange: &iaas.VPCNetworkRange{},
			},
			wantErr: true,
		},
		{
			name: "empty ipv4 network range",
			args: args{
				networkRange: &iaas.VPCNetworkRange{
					VPCNetworkRangeIPv4: &iaas.VPCNetworkRangeIPv4{},
				},
			},
			wantErr: true,
		},
		{
			name: "ipv4 network range with ID",
			args: args{
				networkRange: &iaas.VPCNetworkRange{
					VPCNetworkRangeIPv4: &iaas.VPCNetworkRangeIPv4{
						Id: new(testNetworkRangeId),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty ipv6 network range",
			args: args{
				networkRange: &iaas.VPCNetworkRange{
					VPCNetworkRangeIPv6: &iaas.VPCNetworkRangeIPv6{},
				},
			},
			wantErr: true,
		},
		{
			name: "ipv6 network range with ID",
			args: args{
				networkRange: &iaas.VPCNetworkRange{
					VPCNetworkRangeIPv6: &iaas.VPCNetworkRangeIPv6{
						Id: new(testNetworkRangeId),
					},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.vpcLabel, tt.args.networkRange); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
