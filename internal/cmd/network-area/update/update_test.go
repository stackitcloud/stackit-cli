package update

import (
	"context"
	"strconv"
	"strings"
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

const (
	testRegion                    = "eu01"
	testName                      = "example-network-area-name"
	testDefaultPrefixLength int64 = 25
	testMinPrefixLength     int64 = 24
	testMaxPrefixLength     int64 = 26
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var (
	testOrgId  = uuid.NewString()
	testAreaId = uuid.NewString()

	testDnsNameservers = []string{"1.1.1.0", "1.1.2.0"}
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testAreaId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

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
		Name:           utils.Ptr(testName),
		OrganizationId: utils.Ptr(testOrgId),
		AreaId:         testAreaId,
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiPartialUpdateNetworkAreaRequest)) iaas.ApiPartialUpdateNetworkAreaRequest {
	request := testClient.PartialUpdateNetworkArea(testCtx, testOrgId, testAreaId)
	request = request.PartialUpdateNetworkAreaPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.PartialUpdateNetworkAreaPayload)) iaas.PartialUpdateNetworkAreaPayload {
	payload := iaas.PartialUpdateNetworkAreaPayload{
		Name: utils.Ptr(testName),
		Labels: utils.Ptr(map[string]interface{}{
			"key": "value",
		}),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureRequestRegionalArea(mods ...func(request *iaas.ApiUpdateNetworkAreaRegionRequest)) iaas.ApiUpdateNetworkAreaRegionRequest {
	request := testClient.UpdateNetworkAreaRegion(testCtx, testOrgId, testAreaId, testRegion)
	request = request.UpdateNetworkAreaRegionPayload(fixturePayloadRegionalArea())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayloadRegionalArea(mods ...func(payload *iaas.UpdateNetworkAreaRegionPayload)) iaas.UpdateNetworkAreaRegionPayload {
	payload := iaas.UpdateNetworkAreaRegionPayload{
		Ipv4: &iaas.UpdateRegionalAreaIPv4{
			DefaultNameservers: utils.Ptr(testDnsNameservers),
			DefaultPrefixLen:   utils.Ptr(testDefaultPrefixLength),
			MaxPrefixLen:       utils.Ptr(testMaxPrefixLength),
			MinPrefixLen:       utils.Ptr(testMinPrefixLength),
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
			description: "with deprecated flags",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[dnsNameServersFlag] = strings.Join(testDnsNameservers, ",")
				flagValues[defaultPrefixLengthFlag] = strconv.FormatInt(testDefaultPrefixLength, 10)
				flagValues[maxPrefixLengthFlag] = strconv.FormatInt(testMaxPrefixLength, 10)
				flagValues[minPrefixLengthFlag] = strconv.FormatInt(testMinPrefixLength, 10)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DnsNameServers = utils.Ptr(testDnsNameservers)
				model.DefaultPrefixLength = utils.Ptr(testDefaultPrefixLength)
				model.MaxPrefixLength = utils.Ptr(testMaxPrefixLength)
				model.MinPrefixLength = utils.Ptr(testMinPrefixLength)
			}),
		},

		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "org id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "org id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "org id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "area id missing",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "area id invalid 1",
			argValues: fixtureArgValues(func(argValues []string) {
				argValues[0] = ""
			}),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[areaIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "area id invalid 2",
			argValues: fixtureArgValues(func(argValues []string) {
				argValues[0] = "invalid-uuid"
			}),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[areaIdFlag] = "invalid-uuid"
			}),
			isValid: false,
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
		expectedRequest iaas.ApiPartialUpdateNetworkAreaRequest
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
		expectedRequest iaas.ApiUpdateNetworkAreaRegionRequest
	}{
		{
			description: "base",
			model: fixtureInputModel(func(model *inputModel) {
				model.DnsNameServers = utils.Ptr(testDnsNameservers)
				model.DefaultPrefixLength = utils.Ptr(testDefaultPrefixLength)
				model.MaxPrefixLength = utils.Ptr(testMaxPrefixLength)
				model.MinPrefixLength = utils.Ptr(testMinPrefixLength)
			}),
			expectedRequest: fixtureRequestRegionalArea(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequestNetworkAreaRegion(testCtx, tt.model, testClient)

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

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		responses    NetworkAreaResponses
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
			name: "empty network area",
			args: args{
				responses: NetworkAreaResponses{
					NetworkArea:  iaas.NetworkArea{},
					RegionalArea: nil,
				},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectLabel, tt.args.responses); (err != nil) != tt.wantErr {
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
					OrganizationId: utils.Ptr(testOrgId),
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      nil,
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
					OrganizationId: utils.Ptr(testOrgId),
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      utils.Ptr(testDnsNameservers),
					DefaultPrefixLength: utils.Ptr(testDefaultPrefixLength),
					MaxPrefixLength:     utils.Ptr(testMaxPrefixLength),
					MinPrefixLength:     utils.Ptr(testMinPrefixLength),
				},
			},
			want: []string{dnsNameServersFlag, defaultPrefixLengthFlag, minPrefixLengthFlag, maxPrefixLengthFlag},
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
					OrganizationId: utils.Ptr(testOrgId),
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      nil,
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
					OrganizationId: utils.Ptr(testOrgId),
					Labels: utils.Ptr(map[string]string{
						"key": "value",
					}),
					DnsNameServers:      utils.Ptr(testDnsNameservers),
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
