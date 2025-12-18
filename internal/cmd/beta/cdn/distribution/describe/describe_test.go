package describe

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

type testCtxKey struct{}

var (
	testCtx            = context.WithValue(context.Background(), testCtxKey{}, "test")
	testProjectID      = uuid.NewString()
	testDistributionID = uuid.NewString()
	testClient         = &cdn.APIClient{}
	testTime           = time.Time{}
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectID,
			Verbosity: globalflags.VerbosityDefault,
		},
		DistributionID: testDistributionID,
		WithWAF:        false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureResponse(mods ...func(resp *cdn.GetDistributionResponse)) *cdn.GetDistributionResponse {
	response := &cdn.GetDistributionResponse{
		Distribution: &cdn.Distribution{
			Config: &cdn.Config{
				Backend: &cdn.ConfigBackend{
					BucketBackend: &cdn.BucketBackend{
						BucketUrl: utils.Ptr("https://example.com"),
						Region:    utils.Ptr("eu"),
						Type:      utils.Ptr("bucket"),
					},
				},
				BlockedCountries:     utils.Ptr([]string{}),
				BlockedIps:           utils.Ptr([]string{}),
				DefaultCacheDuration: nil,
				LogSink:              nil,
				MonthlyLimitBytes:    nil,
				Optimizer:            nil,
				Regions:              &[]cdn.Region{cdn.REGION_EU},
				Waf:                  nil,
			},
			CreatedAt: utils.Ptr(testTime),
			Domains:   &[]cdn.Domain{},
			Errors:    nil,
			Id:        utils.Ptr(testDistributionID),
			ProjectId: utils.Ptr(testProjectID),
			Status:    utils.Ptr(cdn.DISTRIBUTIONSTATUS_ACTIVE),
			UpdatedAt: utils.Ptr(testTime),
			Waf:       nil,
		},
	}
	for _, mod := range mods {
		mod(response)
	}
	return response
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		flags       map[string]string
		isValid     bool
		expected    *inputModel
	}{
		{
			description: "base",
			args:        []string{testDistributionID},
			flags:       fixtureFlagValues(),
			isValid:     true,
			expected:    fixtureInputModel(),
		},
		{
			description: "no args",
			args:        []string{},
			flags:       fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid distribution id",
			args:        []string{"invalid-uuid"},
			flags:       fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "missing project id",
			args:        []string{testDistributionID},
			flags:       map[string]string{},
			isValid:     false,
		},
		{
			description: "invalid project id",
			args:        []string{testDistributionID},
			flags: map[string]string{
				globalflags.ProjectIdFlag: "invalid-uuid",
			},
			isValid: false,
		},
		{
			description: "with WAF",
			args:        []string{testDistributionID},
			flags: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[flagWithWaf] = "true"
			}),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.WithWAF = true
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expected, tt.args, tt.flags, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expected    cdn.ApiGetDistributionRequest
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			expected:    testClient.GetDistribution(testCtx, testProjectID, testDistributionID).WithWafStatus(false),
		},
		{
			description: "with WAF",
			model: fixtureInputModel(func(model *inputModel) {
				model.WithWAF = true
			}),
			expected: testClient.GetDistribution(testCtx, testProjectID, testDistributionID).WithWafStatus(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(got, tt.expected,
				cmp.AllowUnexported(tt.expected),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description  string
		format       string
		distribution *cdn.GetDistributionResponse
		wantErr      bool
		expected     string
	}{
		{
			description: "empty",
			format:      "table",
			wantErr:     true,
		},
		{
			description:  "no errors",
			format:       "table",
			distribution: fixtureResponse(),
			//nolint:staticcheck //you can't use escape sequences in ``-string-literals
			expected: fmt.Sprintf(`
[106;30;1m                          Distribution                         [0m
 ID                     │ %-37s
────────────────────────┼──────────────────────────────────────
 STATUS                 │ ACTIVE                               
────────────────────────┼──────────────────────────────────────
 REGIONS                │ EU                                   
────────────────────────┼──────────────────────────────────────
 CREATED AT             │ %-37s
────────────────────────┼──────────────────────────────────────
 UPDATED AT             │ %-37s
────────────────────────┼──────────────────────────────────────
 PROJECT ID             │ %-37s
────────────────────────┼──────────────────────────────────────
 BACKEND TYPE           │ BUCKET                               
────────────────────────┼──────────────────────────────────────
 BUCKET URL             │ https://example.com                  
────────────────────────┼──────────────────────────────────────
 BUCKET REGION          │ eu                                   
────────────────────────┼──────────────────────────────────────
 BLOCKED COUNTRIES      │                                      
────────────────────────┼──────────────────────────────────────
 BLOCKED IPS            │                                      
────────────────────────┼──────────────────────────────────────
 DEFAULT CACHE DURATION │                                      
────────────────────────┼──────────────────────────────────────
 LOG SINK PUSH URL      │                                      
────────────────────────┼──────────────────────────────────────
 MONTHLY LIMIT (BYTES)  │                                      
────────────────────────┼──────────────────────────────────────
 OPTIMIZER ENABLED      │                                      

`,
				testDistributionID,
				testTime,
				testTime,
				testProjectID),
		},
		{
			description: "with errors",
			format:      "table",
			distribution: fixtureResponse(
				func(r *cdn.GetDistributionResponse) {
					r.Distribution.Errors = &[]cdn.StatusError{
						{
							En: utils.Ptr("First error message"),
						},
						{
							En: utils.Ptr("Second error message"),
						},
					}
				},
			),
			//nolint:staticcheck //you can't use escape sequences in ``-string-literals
			expected: fmt.Sprintf(`
[106;30;1m                          Distribution                         [0m
 ID                     │ %-37s
────────────────────────┼──────────────────────────────────────
 STATUS                 │ ACTIVE                               
────────────────────────┼──────────────────────────────────────
 REGIONS                │ EU                                   
────────────────────────┼──────────────────────────────────────
 CREATED AT             │ %-37s
────────────────────────┼──────────────────────────────────────
 UPDATED AT             │ %-37s
────────────────────────┼──────────────────────────────────────
 PROJECT ID             │ %-37s
────────────────────────┼──────────────────────────────────────
 ERRORS                 │ First error message                  
                        │ Second error message                 
────────────────────────┼──────────────────────────────────────
 BACKEND TYPE           │ BUCKET                               
────────────────────────┼──────────────────────────────────────
 BUCKET URL             │ https://example.com                  
────────────────────────┼──────────────────────────────────────
 BUCKET REGION          │ eu                                   
────────────────────────┼──────────────────────────────────────
 BLOCKED COUNTRIES      │                                      
────────────────────────┼──────────────────────────────────────
 BLOCKED IPS            │                                      
────────────────────────┼──────────────────────────────────────
 DEFAULT CACHE DURATION │                                      
────────────────────────┼──────────────────────────────────────
 LOG SINK PUSH URL      │                                      
────────────────────────┼──────────────────────────────────────
 MONTHLY LIMIT (BYTES)  │                                      
────────────────────────┼──────────────────────────────────────
 OPTIMIZER ENABLED      │                                      

`, testDistributionID,
				testTime,
				testTime,
				testProjectID),
		},
		{
			description: "full",
			format:      "table",
			distribution: fixtureResponse(
				func(r *cdn.GetDistributionResponse) {
					r.Distribution.Waf = &cdn.DistributionWaf{
						EnabledRules: &[]cdn.WafStatusRuleBlock{
							{Id: utils.Ptr("rule-id-1")},
							{Id: utils.Ptr("rule-id-2")},
						},
						DisabledRules: &[]cdn.WafStatusRuleBlock{
							{Id: utils.Ptr("rule-id-3")},
							{Id: utils.Ptr("rule-id-4")},
						},
						LogOnlyRules: &[]cdn.WafStatusRuleBlock{
							{Id: utils.Ptr("rule-id-5")},
							{Id: utils.Ptr("rule-id-6")},
						},
					}
					r.Distribution.Config.Backend = &cdn.ConfigBackend{
						HttpBackend: &cdn.HttpBackend{
							OriginUrl: utils.Ptr("https://origin.example.com"),
							OriginRequestHeaders: &map[string]string{
								"X-Custom-Header": "CustomValue",
							},
							Geofencing: &map[string][]string{
								"origin1.example.com": {"US", "CA"},
								"origin2.example.com": {"FR", "DE"},
							},
						},
					}
					r.Distribution.Config.BlockedCountries = &[]string{"US", "CN"}
					r.Distribution.Config.BlockedIps = &[]string{"127.0.0.1"}
					r.Distribution.Config.DefaultCacheDuration = cdn.NewNullableString(utils.Ptr("P1DT2H30M"))
					r.Distribution.Config.LogSink = &cdn.ConfigLogSink{
						LokiLogSink: &cdn.LokiLogSink{
							PushUrl: utils.Ptr("https://logs.example.com"),
						},
					}
					r.Distribution.Config.MonthlyLimitBytes = utils.Ptr(int64(104857600))
					r.Distribution.Config.Optimizer = &cdn.Optimizer{
						Enabled: utils.Ptr(true),
					}
				}),
			//nolint:staticcheck //you can't use escape sequences in ``-string-literals
			expected: fmt.Sprintf(`
[106;30;1m                            Distribution                            [0m
 ID                          │ %-37s
─────────────────────────────┼──────────────────────────────────────
 STATUS                      │ ACTIVE                               
─────────────────────────────┼──────────────────────────────────────
 REGIONS                     │ EU                                   
─────────────────────────────┼──────────────────────────────────────
 CREATED AT                  │ %-37s
─────────────────────────────┼──────────────────────────────────────
 UPDATED AT                  │ %-37s
─────────────────────────────┼──────────────────────────────────────
 PROJECT ID                  │ %-37s
─────────────────────────────┼──────────────────────────────────────
 BACKEND TYPE                │ HTTP                                 
─────────────────────────────┼──────────────────────────────────────
 HTTP ORIGIN URL             │ https://origin.example.com           
─────────────────────────────┼──────────────────────────────────────
 HTTP ORIGIN REQUEST HEADERS │ X-Custom-Header: CustomValue         
─────────────────────────────┼──────────────────────────────────────
 HTTP GEOFENCING PROPERTIES  │ origin1.example.com: US, CA          
                             │ origin2.example.com: FR, DE          
─────────────────────────────┼──────────────────────────────────────
 BLOCKED COUNTRIES           │ US, CN                               
─────────────────────────────┼──────────────────────────────────────
 BLOCKED IPS                 │ 127.0.0.1                            
─────────────────────────────┼──────────────────────────────────────
 DEFAULT CACHE DURATION      │ P1DT2H30M                            
─────────────────────────────┼──────────────────────────────────────
 LOG SINK PUSH URL           │ https://logs.example.com             
─────────────────────────────┼──────────────────────────────────────
 MONTHLY LIMIT (BYTES)       │ 104857600                            
─────────────────────────────┼──────────────────────────────────────
 OPTIMIZER ENABLED           │ true                                 


[106;30;1m              WAF             [0m
 DISABLED RULE ID │ rule-id-3 
──────────────────┼───────────
 DISABLED RULE ID │ rule-id-4 
──────────────────┼───────────
 ENABLED RULE ID  │ rule-id-1 
──────────────────┼───────────
 ENABLED RULE ID  │ rule-id-2 
──────────────────┼───────────
 LOG-ONLY RULE ID │ rule-id-5 
──────────────────┼───────────
 LOG-ONLY RULE ID │ rule-id-6 

`, testDistributionID, testTime, testTime, testProjectID),
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			p.Cmd.SetOut(&buf)
			if err := outputResult(p, tt.format, tt.distribution); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			diff := cmp.Diff(buf.String(), tt.expected)
			if diff != "" {
				t.Fatalf("outputResult() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
