package create

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
	"k8s.io/utils/ptr"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &cdn.APIClient{}
var testProjectId = uuid.NewString()

const testRegions = cdn.REGION_EU

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		flagRegion:                string(testRegions),
	}
	flagsHTTPBackend()(flagValues)
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func flagsHTTPBackend() func(flagValues map[string]string) {
	return func(flagValues map[string]string) {
		delete(flagValues, flagBucket)
		flagValues[flagHTTP] = "true"
		flagValues[flagHTTPOriginURL] = "https://http-backend.example.com"
	}
}

func flagsBucketBackend() func(flagValues map[string]string) {
	return func(flagValues map[string]string) {
		delete(flagValues, flagHTTP)
		flagValues[flagBucket] = "true"
		flagValues[flagBucketURL] = "https://bucket-backend.example.com"
		flagValues[flagBucketCredentialsAccessKeyID] = "access-key-id"
		flagValues[flagBucketRegion] = "eu"
	}
}

func flagsLoki() func(flagValues map[string]string) {
	return func(flagValues map[string]string) {
		flagValues[flagLoki] = "true"
		flagValues[flagLokiPushURL] = "https://loki.example.com"
		flagValues[flagLokiUsername] = "loki-user"
	}
}

func flagRegions(regions ...cdn.Region) func(flagValues map[string]string) {
	return func(flagValues map[string]string) {
		if len(regions) == 0 {
			delete(flagValues, flagRegion)
			return
		}
		stringRegions := sdkUtils.EnumSliceToStringSlice(regions)
		flagValues[flagRegion] = strings.Join(stringRegions, ",")
	}
}

func fixtureModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		Regions: []cdn.Region{testRegions},
	}
	modelHTTPBackend()(model)
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func modelRegions(regions ...cdn.Region) func(model *inputModel) {
	return func(model *inputModel) {
		model.Regions = regions
	}
}

func modelHTTPBackend() func(model *inputModel) {
	return func(model *inputModel) {
		model.Bucket = nil
		model.HTTP = &httpInputModel{
			OriginURL: "https://http-backend.example.com",
		}
	}
}

func modelBucketBackend() func(model *inputModel) {
	return func(model *inputModel) {
		model.HTTP = nil
		model.Bucket = &bucketInputModel{
			URL:         "https://bucket-backend.example.com",
			AccessKeyID: "access-key-id",
			Region:      "eu",
		}
	}
}

func modelLoki() func(model *inputModel) {
	return func(model *inputModel) {
		model.Loki = &lokiInputModel{
			PushURL:  "https://loki.example.com",
			Username: "loki-user",
		}
	}
}

func fixturePayload(mods ...func(payload *cdn.CreateDistributionPayload)) cdn.CreateDistributionPayload {
	payload := *cdn.NewCreateDistributionPayload(
		cdn.CreateDistributionPayloadGetBackendArgType{
			HttpBackendCreate: &cdn.HttpBackendCreate{
				Type:      utils.Ptr("http"),
				OriginUrl: utils.Ptr("https://http-backend.example.com"),
			},
		},
		[]cdn.Region{testRegions},
	)
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func payloadRegions(regions ...cdn.Region) func(payload *cdn.CreateDistributionPayload) {
	return func(payload *cdn.CreateDistributionPayload) {
		payload.Regions = &regions
	}
}

func payloadBucketBackend() func(payload *cdn.CreateDistributionPayload) {
	return func(payload *cdn.CreateDistributionPayload) {
		payload.Backend = &cdn.CreateDistributionPayloadGetBackendArgType{
			BucketBackendCreate: &cdn.BucketBackendCreate{
				Type:      utils.Ptr("bucket"),
				BucketUrl: utils.Ptr("https://bucket-backend.example.com"),
				Region:    utils.Ptr("eu"),
				Credentials: cdn.NewBucketCredentials(
					"access-key-id",
					"",
				),
			},
		}
	}
}

func payloadLoki() func(payload *cdn.CreateDistributionPayload) {
	return func(payload *cdn.CreateDistributionPayload) {
		payload.LogSink = &cdn.CreateDistributionPayloadGetLogSinkArgType{
			LokiLogSinkCreate: &cdn.LokiLogSinkCreate{
				Type:        utils.Ptr("loki"),
				PushUrl:     utils.Ptr("https://loki.example.com"),
				Credentials: cdn.NewLokiLogSinkCredentials("", "loki-user"),
			},
		}
	}
}

func fixtureRequest(mods ...func(payload *cdn.CreateDistributionPayload)) cdn.ApiCreateDistributionRequest {
	req := testClient.CreateDistribution(testCtx, testProjectId)
	req = req.CreateDistributionPayload(fixturePayload(mods...))
	return req
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description string
		argValues   []string
		flagValues  map[string]string
		isValid     bool
		expected    *inputModel
	}{
		{
			description: "base",
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expected:    fixtureModel(),
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
			description: "regions missing",
			flagValues:  fixtureFlagValues(flagRegions()),
			isValid:     false,
		},
		{
			description: "multiple regions",
			flagValues:  fixtureFlagValues(flagRegions(cdn.REGION_EU, cdn.REGION_AF)),
			isValid:     true,
			expected:    fixtureModel(modelRegions(cdn.REGION_EU, cdn.REGION_AF)),
		},
		{
			description: "bucket backend",
			flagValues:  fixtureFlagValues(flagsBucketBackend()),
			isValid:     true,
			expected:    fixtureModel(modelBucketBackend()),
		},
		{
			description: "bucket backend missing url",
			flagValues: fixtureFlagValues(
				flagsBucketBackend(),
				func(flagValues map[string]string) {
					delete(flagValues, flagBucketURL)
				},
			),
			isValid: false,
		},
		{
			description: "bucket backend missing access key id",
			flagValues: fixtureFlagValues(
				flagsBucketBackend(),
				func(flagValues map[string]string) {
					delete(flagValues, flagBucketCredentialsAccessKeyID)
				},
			),
			isValid: false,
		},
		{
			description: "bucket backend missing region",
			flagValues: fixtureFlagValues(
				flagsBucketBackend(),
				func(flagValues map[string]string) {
					delete(flagValues, flagBucketRegion)
				},
			),
			isValid: false,
		},
		{
			description: "http backend missing url",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, flagHTTPOriginURL)
				},
			),
			isValid: false,
		},
		{
			description: "http backend with geofencing",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagHTTPGeofencing] = "https://dach.example.com DE,AT,CH"
				},
			),
			isValid: true,
			expected: fixtureModel(
				func(model *inputModel) {
					model.HTTP.Geofencing = &map[string][]string{
						"https://dach.example.com": {"DE", "AT", "CH"},
					}
				},
			),
		},
		{
			description: "http backend with origin request headers",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagHTTPOriginRequestHeaders] = "X-Custom-Header:Value1,X-Another-Header:Value2"
				},
			),
			isValid: true,
			expected: fixtureModel(
				func(model *inputModel) {
					model.HTTP.OriginRequestHeaders = &map[string]string{
						"X-Custom-Header":  "Value1",
						"X-Another-Header": "Value2",
					}
				},
			),
		},
		{
			description: "with blocked countries",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagBlockedCountries] = "DE,AT"
				}),
			isValid: true,
			expected: fixtureModel(
				func(model *inputModel) {
					model.BlockedCountries = []string{"DE", "AT"}
				},
			),
		},
		{
			description: "with blocked IPs",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagBlockedIPs] = "127.0.0.1,10.0.0.8"
				}),
			isValid: true,
			expected: fixtureModel(
				func(model *inputModel) {
					model.BlockedIPs = []string{"127.0.0.1", "10.0.0.8"}
				}),
		},
		{
			description: "with default cache duration",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagDefaultCacheDuration] = "PT1H30M"
				}),
			isValid: true,
			expected: fixtureModel(
				func(model *inputModel) {
					model.DefaultCacheDuration = "PT1H30M"
				}),
		},
		{
			description: "with optimizer",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagOptimizer] = "true"
				}),
			isValid: true,
			expected: fixtureModel(
				func(model *inputModel) {
					model.Optimizer = true
				}),
		},
		{
			description: "with loki",
			flagValues: fixtureFlagValues(
				flagsLoki(),
			),
			isValid: true,
			expected: fixtureModel(
				modelLoki(),
			),
		},
		{
			description: "loki with missing username",
			flagValues: fixtureFlagValues(
				flagsLoki(),
				func(flagValues map[string]string) {
					delete(flagValues, flagLokiUsername)
				},
			),
			isValid: false,
		},
		{
			description: "loki with missing push url",
			flagValues: fixtureFlagValues(
				flagsLoki(),
				func(flagValues map[string]string) {
					delete(flagValues, flagLokiPushURL)
				},
			),
			isValid: false,
		},
		{
			description: "with monthly limit bytes",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagMonthlyLimitBytes] = "1073741824" // 1 GiB
				}),
			isValid: true,
			expected: fixtureModel(
				func(model *inputModel) {
					model.MonthlyLimitBytes = ptr.To[int64](1073741824)
				}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expected, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expected    cdn.ApiCreateDistributionRequest
	}{
		{
			description: "base",
			model:       fixtureModel(),
			expected:    fixtureRequest(),
		},
		{
			description: "multiple regions",
			model:       fixtureModel(modelRegions(cdn.REGION_AF, cdn.REGION_EU)),
			expected:    fixtureRequest(payloadRegions(cdn.REGION_AF, cdn.REGION_EU)),
		},
		{
			description: "bucket backend",
			model:       fixtureModel(modelBucketBackend()),
			expected:    fixtureRequest(payloadBucketBackend()),
		},
		{
			description: "http backend with geofencing and origin request headers",
			model: fixtureModel(
				func(model *inputModel) {
					model.HTTP.Geofencing = &map[string][]string{
						"https://dach.example.com": {"DE", "AT", "CH"},
					}
					model.HTTP.OriginRequestHeaders = &map[string]string{
						"X-Custom-Header":  "Value1",
						"X-Another-Header": "Value2",
					}
				},
			),
			expected: fixtureRequest(
				func(payload *cdn.CreateDistributionPayload) {
					payload.Backend.HttpBackendCreate.Geofencing = &map[string][]string{
						"https://dach.example.com": {"DE", "AT", "CH"},
					}
					payload.Backend.HttpBackendCreate.OriginRequestHeaders = &map[string]string{
						"X-Custom-Header":  "Value1",
						"X-Another-Header": "Value2",
					}
				},
			),
		},
		{
			description: "with full options",
			model: fixtureModel(
				func(model *inputModel) {
					model.MonthlyLimitBytes = ptr.To[int64](5368709120) // 5 GiB
					model.Optimizer = true
					model.BlockedCountries = []string{"DE", "AT"}
					model.BlockedIPs = []string{"127.0.0.1"}
					model.DefaultCacheDuration = "PT2H"
				},
			),
			expected: fixtureRequest(
				func(payload *cdn.CreateDistributionPayload) {
					payload.MonthlyLimitBytes = utils.Ptr[int64](5368709120)
					payload.Optimizer = &cdn.CreateDistributionPayloadGetOptimizerArgType{
						Enabled: utils.Ptr(true),
					}
					payload.BlockedCountries = &[]string{"DE", "AT"}
					payload.BlockedIps = &[]string{"127.0.0.1"}
					payload.DefaultCacheDuration = utils.Ptr("PT2H")
				},
			),
		},
		{
			description: "loki",
			model: fixtureModel(
				modelLoki(),
			),
			expected: fixtureRequest(payloadLoki()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expected,
				cmp.AllowUnexported(tt.expected),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description  string
		outputFormat string
		response     *cdn.CreateDistributionResponse
		expected     string
		wantErr      bool
	}{
		{
			description:  "nil response",
			outputFormat: "table",
			response:     nil,
			wantErr:      true,
		},
		{
			description:  "table output",
			outputFormat: "table",
			response: &cdn.CreateDistributionResponse{
				Distribution: &cdn.Distribution{
					Id: ptr.To("dist-1234"),
				},
			},
			expected: fmt.Sprintf("Created CDN distribution for %q. ID: dist-1234\n", testProjectId),
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			p.Cmd.SetOut(buffer)
			if err := outputResult(p, tt.outputFormat, testProjectId, tt.response); (err != nil) != tt.wantErr {
				t.Fatalf("outputResult: %v", err)
			}
			if buffer.String() != tt.expected {
				t.Errorf("want:\n%s\ngot:\n%s", tt.expected, buffer.String())
			}
		})
	}
}
