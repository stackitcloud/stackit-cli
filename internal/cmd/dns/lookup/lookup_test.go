package lookup

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

const (
	testDomain = "www.example.com"
	testType   = dns.LOOKUPDOMAINPAYLOADRECORDTYPE_AAAA
	testRegion = "eu01"
)

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{DefaultAPI: &dns.DefaultAPIService{}}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testDomain,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,
		recordTypeFlag.Name():  string(testType),
		allResolversFlag:       "true",
		includeRawResponseFlag: "true",
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
		Domain:             testDomain,
		RecordType:         utils.Ptr(testType),
		AllResolvers:       true,
		IncludeRawResponse: true,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(payload *dns.LookupDomainPayload)) dns.ApiLookupDomainRequest {
	payload := dns.LookupDomainPayload{
		Address:            utils.Ptr(testDomain),
		RecordType:         utils.Ptr(testType),
		Type:               utils.Ptr(dns.LOOKUPDOMAINPAYLOADTYPE_ALL),
		IncludeRawResponse: utils.Ptr(true),
		Region:             utils.Ptr(testRegion),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return testClient.DefaultAPI.LookupDomain(testCtx).LookupDomainPayload(payload)
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "only required values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Region = ""
				model.RecordType = nil
				model.AllResolvers = false
				model.IncludeRawResponse = false
			}),
		},
		{
			description: "include raw response with json output",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.OutputFormatFlag.Name()] = print.JSONOutputFormat
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.OutputFormat = print.JSONOutputFormat
			}),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "domain empty",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "record type invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[recordTypeFlag.Name()] = "INVALID"
			}),
			isValid: false,
		},
		{
			description: "record type empty",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[recordTypeFlag.Name()] = ""
			}),
			isValid: false,
		},
		{
			description: "all resolvers invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[allResolversFlag] = "invalid"
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
		expectedRequest dns.ApiLookupDomainRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "only required values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Region = ""
				model.RecordType = nil
				model.AllResolvers = false
				model.IncludeRawResponse = false
			}),
			expectedRequest: fixtureRequest(func(payload *dns.LookupDomainPayload) {
				payload.Region = nil
				payload.RecordType = nil
				payload.Type = nil
				payload.IncludeRawResponse = nil
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, dns.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestParseInputRawResponseWarning(t *testing.T) {
	tests := []struct {
		description  string
		flagValues   map[string]string
		expectedWarn bool
	}{
		{
			description:  "default output format",
			flagValues:   fixtureFlagValues(),
			expectedWarn: true,
		},
		{
			description: "pretty output format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.OutputFormatFlag.Name()] = print.PrettyOutputFormat
			}),
			expectedWarn: true,
		},
		{
			description: "none output format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.OutputFormatFlag.Name()] = print.NoneOutputFormat
			}),
			expectedWarn: true,
		},
		{
			description: "json output format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.OutputFormatFlag.Name()] = print.JSONOutputFormat
			}),
			expectedWarn: false,
		},
		{
			description: "yaml output format",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.OutputFormatFlag.Name()] = print.YAMLOutputFormat
			}),
			expectedWarn: false,
		},
		{
			description: "without include raw response",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, includeRawResponseFlag)
			}),
			expectedWarn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			cmd := NewCmd(params.CmdParams)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}
			recordTypeFlag.Reset()
			globalflags.OutputFormatFlag.Reset()
			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			_, err = parseInput(params.Printer, cmd, fixtureArgValues())
			if err != nil {
				t.Fatalf("parseInput() error = %v", err)
			}

			warned := strings.Contains(params.Err.String(), includeRawResponseFlag)
			if warned != tt.expectedWarn {
				t.Errorf("expected warning = %t, got %t (stderr: %q)", tt.expectedWarn, warned, params.Err.String())
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	const (
		testResolverIP      = "192.0.2.1"
		testEmptyResolverIP = "192.0.2.2"
		testRecordData      = "2001:db8::1"
		testRawResponse     = ";; opcode: QUERY, status: NOERROR"
	)
	fixtureResponse := func() *dns.LookupDomainResponse {
		return &dns.LookupDomainResponse{
			State: []dns.LookupResult{
				{
					ResolverIP:  utils.Ptr(testResolverIP),
					RawResponse: utils.Ptr(testRawResponse),
					Records: []dns.LookupRecord{
						{
							Name: utils.Ptr(testDomain),
							Type: utils.Ptr(dns.LOOKUPRECORDTYPE_AAAA),
							Ttl:  utils.Ptr(int32(3600)),
							Data: utils.Ptr(testRecordData),
						},
						{},
					},
				},
			},
		}
	}

	type args struct {
		model *inputModel
		resp  *dns.LookupDomainResponse
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantOutput    []string
		notWantOutput []string
	}{
		{
			name: "empty",
			args: args{
				model: fixtureInputModel(),
			},
			wantErr: true,
		},
		{
			name: "no results",
			args: args{
				model: fixtureInputModel(),
				resp:  &dns.LookupDomainResponse{},
			},
			wantOutput: []string{`No records found for domain "www.example.com"`},
		},
		{
			name: "single resolver without records",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.AllResolvers = false
				}),
				resp: &dns.LookupDomainResponse{
					State: []dns.LookupResult{{ResolverIP: utils.Ptr(testEmptyResolverIP)}},
				},
			},
			wantOutput:    []string{`No records found for domain "www.example.com"`},
			notWantOutput: []string{testEmptyResolverIP},
		},
		{
			name: "all resolvers with one resolver without records",
			args: args{
				model: fixtureInputModel(),
				resp: func() *dns.LookupDomainResponse {
					resp := fixtureResponse()
					resp.State = append(resp.State, dns.LookupResult{ResolverIP: utils.Ptr(testEmptyResolverIP)})
					return resp
				}(),
			},
			wantOutput:    []string{testResolverIP, testRecordData, testEmptyResolverIP, "No records found"},
			notWantOutput: []string{"No records found for domain"},
		},
		{
			name: "all resolvers without records",
			args: args{
				model: fixtureInputModel(),
				resp: &dns.LookupDomainResponse{
					State: []dns.LookupResult{
						{ResolverIP: utils.Ptr(testResolverIP)},
						{ResolverIP: utils.Ptr(testEmptyResolverIP)},
					},
				},
			},
			wantOutput:    []string{testResolverIP, testEmptyResolverIP, "No records found"},
			notWantOutput: []string{"No records found for domain"},
		},
		{
			name: "pretty",
			args: args{
				model: fixtureInputModel(),
				resp:  fixtureResponse(),
			},
			wantOutput:    []string{testResolverIP, testDomain, testRecordData},
			notWantOutput: []string{testRawResponse},
		},
		{
			name: "json with raw response",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: fixtureResponse(),
			},
			wantOutput: []string{`"rawResponse": "` + testRawResponse + `"`, `"data": "` + testRecordData + `"`},
		},
		{
			name: "yaml with raw response",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
				}),
				resp: fixtureResponse(),
			},
			wantOutput: []string{"rawResponse: \"" + testRawResponse + "\"", "data: " + testRecordData},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()
			if err := outputResult(params.Printer, tt.args.model, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			output := params.Out.String()
			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("output %q does not contain %q", output, want)
				}
			}
			for _, notWant := range tt.notWantOutput {
				if strings.Contains(output, notWant) {
					t.Errorf("output %q contains %q", output, notWant)
				}
			}
		})
	}
}
