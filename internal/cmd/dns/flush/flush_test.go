package flush

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

const (
	testDomain = "www.example.com"
	testType   = dns.FLUSHDOMAINPAYLOADTYPE_AAAA
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
		typeFlag.Name(): string(testType),
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
		Domain: testDomain,
		Type:   utils.Ptr(testType),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiFlushDomainRequest)) dns.ApiFlushDomainRequest {
	request := testClient.DefaultAPI.FlushDomain(testCtx)
	request = request.FlushDomainPayload(dns.FlushDomainPayload{
		Domain: testDomain,
		Type:   utils.Ptr(testType),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			description: "without type",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, typeFlag.Name())
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Type = nil
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
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Type = nil
			}),
		},
		{
			description: "domain empty",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "type invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag.Name()] = "INVALID"
			}),
			isValid: false,
		},
		{
			description: "type lowercase",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag.Name()] = "aaaa"
			}),
			isValid: false,
		},
		{
			description: "type empty",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag.Name()] = ""
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
		expectedRequest dns.ApiFlushDomainRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "without type",
			model: fixtureInputModel(func(model *inputModel) {
				model.Type = nil
			}),
			expectedRequest: fixtureRequest(func(request *dns.ApiFlushDomainRequest) {
				*request = request.FlushDomainPayload(dns.FlushDomainPayload{
					Domain: testDomain,
				})
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

func TestOutputResult(t *testing.T) {
	type args struct {
		model *inputModel
		resp  *dns.FlushDomainResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				model: fixtureInputModel(),
			},
			wantErr: true,
		},
		{
			name: "base",
			args: args{
				model: fixtureInputModel(),
				resp:  &dns.FlushDomainResponse{},
			},
			wantErr: false,
		},
		{
			name: "without type",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Type = nil
				}),
				resp: &dns.FlushDomainResponse{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
