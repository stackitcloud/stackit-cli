package update

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
	"k8s.io/utils/ptr"
)

const testCacheDuration = "P1DT12H"

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &cdn.APIClient{}
var testProjectId = uuid.NewString()
var testDistributionID = uuid.NewString()

const testMonthlyLimitBytes int64 = 1048576

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
	}
	for _, m := range mods {
		m(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
		},
		DistributionID: testDistributionID,
		Regions:        []cdn.Region{},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(payload *cdn.PatchDistributionPayload)) cdn.ApiPatchDistributionRequest {
	req := testClient.PatchDistribution(testCtx, testProjectId, testDistributionID)
	if payload := fixturePayload(mods...); payload != nil {
		req = req.PatchDistributionPayload(*fixturePayload(mods...))
	}
	return req
}

func fixturePayload(mods ...func(payload *cdn.PatchDistributionPayload)) *cdn.PatchDistributionPayload {
	payload := cdn.NewPatchDistributionPayload()
	payload.Config = &cdn.ConfigPatch{}
	for _, m := range mods {
		m(payload)
	}
	return payload
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
			argValues:   []string{testDistributionID},
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expected:    fixtureInputModel(),
		},
		{
			description: "distribution id missing",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid distribution id",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   []string{testDistributionID},
			flagValues:  fixtureFlagValues(func(flagValues map[string]string) { delete(flagValues, globalflags.ProjectIdFlag) }),
			isValid:     false,
		},
		{
			description: "invalid distribution id",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "both backends",
			argValues:   []string{testDistributionID},
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagHTTP] = "true"
					flagValues[flagBucket] = "true"
				},
			),
			isValid: false,
		},
		{
			description: "max config without backend",
			argValues:   []string{testDistributionID},
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagRegions] = "EU,US"
					flagValues[flagBlockedCountries] = "DE,AT,CH"
					flagValues[flagBlockedIPs] = "127.0.0.1,10.0.0.8"
					flagValues[flagDefaultCacheDuration] = "P1DT12H"
					flagValues[flagLoki] = "true"
					flagValues[flagLokiUsername] = "loki-user"
					flagValues[flagLokiPushURL] = "https://loki.example.com"
					flagValues[flagMonthlyLimitBytes] = fmt.Sprintf("%d", testMonthlyLimitBytes)
					flagValues[flagOptimizer] = "true"
				},
			),
			isValid: true,
			expected: fixtureInputModel(
				func(model *inputModel) {
					model.Regions = []cdn.Region{cdn.REGION_EU, cdn.REGION_US}
					model.BlockedCountries = []string{"DE", "AT", "CH"}
					model.BlockedIPs = []string{"127.0.0.1", "10.0.0.8"}
					model.DefaultCacheDuration = "P1DT12H"
					model.Loki = &lokiInputModel{
						Username: "loki-user",
						PushURL:  "https://loki.example.com",
					}
					model.MonthlyLimitBytes = utils.Ptr(testMonthlyLimitBytes)
					model.Optimizer = utils.Ptr(true)
				},
			),
		},
		{
			description: "max config http backend",
			argValues:   []string{testDistributionID},
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagHTTP] = "true"
					flagValues[flagHTTPOriginURL] = "https://origin.example.com"
					flagValues[flagHTTPOriginRequestHeaders] = "X-Example-Header: example-value, X-Another-Header: another-value"
					flagValues[flagHTTPGeofencing] = "https://dach.example.com DE,AT,CH"
				},
			),
			isValid: true,
			expected: fixtureInputModel(
				func(model *inputModel) {
					model.HTTP = &httpInputModel{
						OriginURL: "https://origin.example.com",
						OriginRequestHeaders: &map[string]string{
							"X-Example-Header": "example-value",
							"X-Another-Header": "another-value",
						},
						Geofencing: &map[string][]string{
							"https://dach.example.com": {"DE", "AT", "CH"},
						},
					}
				},
			),
		},
		{
			description: "max config bucket backend",
			argValues:   []string{testDistributionID},
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[flagBucket] = "true"
					flagValues[flagBucketURL] = "https://bucket.example.com"
					flagValues[flagBucketRegion] = "EU"
					flagValues[flagBucketCredentialsAccessKeyID] = "access-key-id"
				},
			),
			isValid: true,
			expected: fixtureInputModel(
				func(model *inputModel) {
					model.Bucket = &bucketInputModel{
						URL:         "https://bucket.example.com",
						Region:      "EU",
						AccessKeyID: "access-key-id",
					}
				},
			),
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
		expected    cdn.ApiPatchDistributionRequest
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			expected:    fixtureRequest(),
		},
		{
			description: "max without backend",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.Regions = []cdn.Region{cdn.REGION_EU, cdn.REGION_US}
					model.BlockedCountries = []string{"DE", "AT", "CH"}
					model.BlockedIPs = []string{"127.0.0.1", "10.0.0.8"}
					model.DefaultCacheDuration = testCacheDuration
					model.MonthlyLimitBytes = utils.Ptr(testMonthlyLimitBytes)
					model.Loki = &lokiInputModel{
						Password: "loki-pass",
						Username: "loki-user",
						PushURL:  "https://loki.example.com",
					}
					model.Optimizer = utils.Ptr(true)
				},
			),
			expected: fixtureRequest(
				func(payload *cdn.PatchDistributionPayload) {
					payload.Config.Regions = &[]cdn.Region{cdn.REGION_EU, cdn.REGION_US}
					payload.Config.BlockedCountries = &[]string{"DE", "AT", "CH"}
					payload.Config.BlockedIps = &[]string{"127.0.0.1", "10.0.0.8"}
					payload.Config.DefaultCacheDuration = cdn.NewNullableString(utils.Ptr(testCacheDuration))
					payload.Config.MonthlyLimitBytes = utils.Ptr(testMonthlyLimitBytes)
					payload.Config.LogSink = cdn.NewNullableConfigPatchLogSink(&cdn.ConfigPatchLogSink{
						LokiLogSinkPatch: &cdn.LokiLogSinkPatch{
							Credentials: cdn.NewLokiLogSinkCredentials("loki-pass", "loki-user"),
							PushUrl:     utils.Ptr("https://loki.example.com"),
						},
					})
					payload.Config.Optimizer = &cdn.OptimizerPatch{
						Enabled: utils.Ptr(true),
					}
				},
			),
		},
		{
			description: "max http backend",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.HTTP = &httpInputModel{
						Geofencing:           &map[string][]string{"https://dach.example.com": {"DE", "AT", "CH"}},
						OriginRequestHeaders: &map[string]string{"X-Example-Header": "example-value", "X-Another-Header": "another-value"},
						OriginURL:            "https://http-backend.example.com",
					}
				}),
			expected: fixtureRequest(
				func(payload *cdn.PatchDistributionPayload) {
					payload.Config.Backend = &cdn.ConfigPatchBackend{
						HttpBackendPatch: &cdn.HttpBackendPatch{
							Geofencing: &map[string][]string{"https://dach.example.com": {"DE", "AT", "CH"}},
							OriginRequestHeaders: &map[string]string{
								"X-Example-Header": "example-value",
								"X-Another-Header": "another-value",
							},
							OriginUrl: utils.Ptr("https://http-backend.example.com"),
							Type:      utils.Ptr("http"),
						},
					}
				}),
		},
		{
			description: "max bucket backend",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.Bucket = &bucketInputModel{
						URL:         "https://bucket.example.com",
						AccessKeyID: "bucket-access-key-id",
						Password:    "bucket-pass",
						Region:      "EU",
					}
				}),
			expected: fixtureRequest(
				func(payload *cdn.PatchDistributionPayload) {
					payload.Config.Backend = &cdn.ConfigPatchBackend{
						BucketBackendPatch: &cdn.BucketBackendPatch{
							BucketUrl:   utils.Ptr("https://bucket.example.com"),
							Credentials: cdn.NewBucketCredentials("bucket-access-key-id", "bucket-pass"),
							Region:      utils.Ptr("EU"),
							Type:        utils.Ptr("bucket"),
						},
					}
				}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, testClient, tt.model)

			diff := cmp.Diff(request, tt.expected,
				cmp.AllowUnexported(tt.expected, cdn.NullableString{}, cdn.NullableConfigPatchLogSink{}),
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
		response     *cdn.PatchDistributionResponse
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
			response: &cdn.PatchDistributionResponse{
				Distribution: &cdn.Distribution{
					Id: ptr.To("dist-1234"),
				},
			},
			expected: fmt.Sprintf("Updated CDN distribution for %q. ID: dist-1234\n", testProjectId),
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
