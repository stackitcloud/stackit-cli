package update

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

var projectIdFlag = globalflags.ProjectIdFlag.FlagName()

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{}
var testProjectId = uuid.NewString()
var testZoneId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:     testProjectId,
		zoneIdFlag:        testZoneId,
		nameFlag:          "example",
		defaultTTLFlag:    "3600",
		aclFlag:           "0.0.0.0/0",
		primaryFlag:       "1.1.1.1",
		retryTimeFlag:     "600",
		refreshTimeFlag:   "3600",
		negativeCacheFlag: "60",
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
		ProjectId:     testProjectId,
		ZoneId:        testZoneId,
		Name:          utils.Ptr("example"),
		DefaultTTL:    utils.Ptr(int64(3600)),
		Primaries:     utils.Ptr([]string{"1.1.1.1"}),
		Acl:           utils.Ptr("0.0.0.0/0"),
		RetryTime:     utils.Ptr(int64(600)),
		RefreshTime:   utils.Ptr(int64(3600)),
		NegativeCache: utils.Ptr(int64(60)),
		ExpireTime:    utils.Ptr(int64(36000000)),
		Description:   utils.Ptr("Example"),
		ContactEmail:  utils.Ptr("example@example.com"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiUpdateZoneRequest)) dns.ApiUpdateZoneRequest {
	request := testClient.UpdateZone(testCtx, testProjectId, testZoneId)
	request = request.UpdateZonePayload(dns.UpdateZonePayload{
		Name:          utils.Ptr("example"),
		DefaultTTL:    utils.Ptr(int64(3600)),
		Primaries:     utils.Ptr([]string{"1.1.1.1"}),
		Acl:           utils.Ptr("0.0.0.0/0"),
		RetryTime:     utils.Ptr(int64(600)),
		RefreshTime:   utils.Ptr(int64(3600)),
		NegativeCache: utils.Ptr(int64(60)),
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
			description: "required flags only (no values to update)",
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
				zoneIdFlag:    testZoneId,
			},
			isValid: false,
			expectedModel: &flagModel{
				ProjectId: testProjectId,
				ZoneId:    testZoneId,
			},
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				projectIdFlag:     testProjectId,
				zoneIdFlag:        testZoneId,
				nameFlag:          "",
				defaultTTLFlag:    "0",
				aclFlag:           "",
				primaryFlag:       "",
				retryTimeFlag:     "0",
				refreshTimeFlag:   "0",
				negativeCacheFlag: "0",
				expireTimeFlag:    "0",
				descriptionFlag:   "",
				contactEmailFlag:  "",
			},
			isValid: true,
			expectedModel: &flagModel{
				ProjectId:     testProjectId,
				ZoneId:        testZoneId,
				Name:          utils.Ptr(""),
				DefaultTTL:    utils.Ptr(int64(0)),
				Primaries:     utils.Ptr([]string{}),
				Acl:           utils.Ptr(""),
				RetryTime:     utils.Ptr(int64(0)),
				RefreshTime:   utils.Ptr(int64(0)),
				NegativeCache: utils.Ptr(int64(0)),
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
			description: "zone id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, zoneIdFlag)
			}),
			isValid: false,
		},
		{
			description: "zone id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[zoneIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "zone id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[zoneIdFlag] = "invalid-uuid"
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
			err := globalflags.ConfigureFlags(cmd.Flags())
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
		expectedRequest dns.ApiUpdateZoneRequest
	}{
		{
			description:     "base",
			model:           fixtureFlagModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "required fields only",
			model: &flagModel{
				ProjectId: testProjectId,
				ZoneId:    testZoneId,
			},
			expectedRequest: testClient.UpdateZone(testCtx, testProjectId, testZoneId).
				UpdateZonePayload(dns.UpdateZonePayload{}),
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
