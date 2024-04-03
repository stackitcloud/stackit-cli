package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{}
var testProjectId = uuid.NewString()
var testZoneId = uuid.NewString()

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testZoneId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:     testProjectId,
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

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
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

func fixtureRequest(mods ...func(request *dns.ApiPartialUpdateZoneRequest)) dns.ApiPartialUpdateZoneRequest {
	request := testClient.PartialUpdateZone(testCtx, testProjectId, testZoneId)
	request = request.PartialUpdateZonePayload(dns.PartialUpdateZonePayload{
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

func TestParseInput(t *testing.T) {
	tests := []struct {
		description       string
		argValues         []string
		flagValues        map[string]string
		primaryFlagValues []string
		isValid           bool
		expectedModel     *inputModel
	}{
		{
			description:   "base",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
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
			isValid:     false,
		},
		{
			description: "required flags only (no values to update)",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
			},
			isValid: false,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				ZoneId: testZoneId,
			},
		},
		{
			description: "zero values",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag:     testProjectId,
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
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
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
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "zone id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "zone id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description:       "repeated primary flags",
			argValues:         fixtureArgValues(),
			flagValues:        fixtureFlagValues(),
			primaryFlagValues: []string{"1.2.3.4", "5.6.7.8"},
			isValid:           true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Primaries = utils.Ptr(
					append(*model.Primaries, "1.2.3.4", "5.6.7.8"),
				)
			}),
		},
		{
			description:       "repeated primary flags with list value",
			argValues:         fixtureArgValues(),
			flagValues:        fixtureFlagValues(),
			primaryFlagValues: []string{"1.2.3.4,5.6.7.8"},
			isValid:           true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Primaries = utils.Ptr(
					append(*model.Primaries, "1.2.3.4", "5.6.7.8"),
				)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd()
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

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
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
		expectedRequest dns.ApiPartialUpdateZoneRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "required fields only",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				ZoneId: testZoneId,
			},
			expectedRequest: testClient.PartialUpdateZone(testCtx, testProjectId, testZoneId).
				PartialUpdateZonePayload(dns.PartialUpdateZonePayload{}),
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
