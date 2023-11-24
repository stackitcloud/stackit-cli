package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{}
var testProjectId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:     testProjectId,
		nameFlag:          "example",
		dnsNameFlag:       "example.com",
		defaultTTLFlag:    "3600",
		aclFlag:           "0.0.0.0/0",
		typeFlag:          "master",
		primaryFlag:       "1.1.1.1",
		retryTimeFlag:     "600",
		refreshTimeFlag:   "3600",
		negativeCacheFlag: "60",
		isReverseZoneFlag: "false",
		expireTimeFlag:    "36000000",
		descriptionFlag:   "Example",
		contactEmailFlag:  "example@example.com",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagModel(mods ...func(model *flagModel)) *flagModel {
	model := &flagModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
		},
		Name:          utils.Ptr("example"),
		DnsName:       utils.Ptr("example.com"),
		DefaultTTL:    utils.Ptr(int64(3600)),
		Primaries:     utils.Ptr([]string{"1.1.1.1"}),
		Acl:           utils.Ptr("0.0.0.0/0"),
		Type:          utils.Ptr("master"),
		RetryTime:     utils.Ptr(int64(600)),
		RefreshTime:   utils.Ptr(int64(3600)),
		NegativeCache: utils.Ptr(int64(60)),
		IsReverseZone: utils.Ptr(false),
		ExpireTime:    utils.Ptr(int64(36000000)),
		Description:   utils.Ptr("Example"),
		ContactEmail:  utils.Ptr("example@example.com"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiCreateZoneRequest)) dns.ApiCreateZoneRequest {
	request := testClient.CreateZone(testCtx, testProjectId)
	request = request.CreateZonePayload(dns.CreateZonePayload{
		Name:          utils.Ptr("example"),
		DnsName:       utils.Ptr("example.com"),
		DefaultTTL:    utils.Ptr(int64(3600)),
		Primaries:     utils.Ptr([]string{"1.1.1.1"}),
		Acl:           utils.Ptr("0.0.0.0/0"),
		Type:          utils.Ptr("master"),
		RetryTime:     utils.Ptr(int64(600)),
		RefreshTime:   utils.Ptr(int64(3600)),
		NegativeCache: utils.Ptr(int64(60)),
		IsReverseZone: utils.Ptr(false),
		ExpireTime:    utils.Ptr(int64(36000000)),
		Description:   utils.Ptr("Example"),
		ContactEmail:  utils.Ptr("example@example.com"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		description       string
		flagValues        map[string]string
		primaryFlagValues []string
		isValid           bool
		expectedModel     *flagModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureFlagModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
				nameFlag:      "example",
				dnsNameFlag:   "example.com",
			},
			isValid: true,
			expectedModel: &flagModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				Name:    utils.Ptr("example"),
				DnsName: utils.Ptr("example.com"),
			},
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				projectIdFlag:     testProjectId,
				nameFlag:          "",
				dnsNameFlag:       "",
				defaultTTLFlag:    "0",
				aclFlag:           "",
				typeFlag:          "",
				primaryFlag:       "",
				retryTimeFlag:     "0",
				refreshTimeFlag:   "0",
				negativeCacheFlag: "0",
				isReverseZoneFlag: "false",
				expireTimeFlag:    "0",
				descriptionFlag:   "",
				contactEmailFlag:  "",
			},
			isValid: true,
			expectedModel: &flagModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				Name:          utils.Ptr(""),
				DnsName:       utils.Ptr(""),
				DefaultTTL:    utils.Ptr(int64(0)),
				Primaries:     utils.Ptr([]string{}),
				Acl:           utils.Ptr(""),
				Type:          utils.Ptr(""),
				RetryTime:     utils.Ptr(int64(0)),
				RefreshTime:   utils.Ptr(int64(0)),
				NegativeCache: utils.Ptr(int64(0)),
				IsReverseZone: utils.Ptr(false),
				ExpireTime:    utils.Ptr(int64(0)),
				Description:   utils.Ptr(""),
				ContactEmail:  utils.Ptr(""),
			},
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description:       "repeated primary flags",
			flagValues:        fixtureFlagValues(),
			primaryFlagValues: []string{"1.2.3.4", "5.6.7.8"},
			isValid:           true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.Primaries = utils.Ptr(
					append(*model.Primaries, "1.2.3.4", "5.6.7.8"),
				)
			}),
		},
		{
			description:       "repeated primary flags with list value",
			flagValues:        fixtureFlagValues(),
			primaryFlagValues: []string{"1.2.3.4,5.6.7.8"},
			isValid:           true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.Primaries = utils.Ptr(
					append(*model.Primaries, "1.2.3.4", "5.6.7.8"),
				)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			configureFlags(cmd)
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

			for _, value := range tt.primaryFlagValues {
				err := cmd.Flags().Set(primaryFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", primaryFlag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseFlags(cmd)
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
		model           *flagModel
		expectedRequest dns.ApiCreateZoneRequest
	}{
		{
			description:     "base",
			model:           fixtureFlagModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "required fields only",
			model: &flagModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				Name:    utils.Ptr("example"),
				DnsName: utils.Ptr("example.com"),
			},
			expectedRequest: testClient.CreateZone(testCtx, testProjectId).
				CreateZonePayload(dns.CreateZonePayload{
					Name:    utils.Ptr("example"),
					DnsName: utils.Ptr("example.com"),
				}),
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
