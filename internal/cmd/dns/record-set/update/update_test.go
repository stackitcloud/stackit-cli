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
var testRecordSetId = uuid.NewString()

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testRecordSetId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		zoneIdFlag:    testZoneId,
		commentFlag:   "comment",
		nameFlag:      "example.com",
		recordFlag:    "1.1.1.1",
		ttlFlag:       "3600",
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
		},
		ZoneId:      testZoneId,
		RecordSetId: testRecordSetId,
		Name:        utils.Ptr("example.com"),
		Comment:     utils.Ptr("comment"),
		Records:     &[]string{"1.1.1.1"},
		TTL:         utils.Ptr(int64(3600)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiPartialUpdateRecordSetRequest)) dns.ApiPartialUpdateRecordSetRequest {
	request := testClient.PartialUpdateRecordSet(testCtx, testProjectId, testZoneId, testRecordSetId)
	request = request.PartialUpdateRecordSetPayload(dns.PartialUpdateRecordSetPayload{
		Name:    utils.Ptr("example.com"),
		Comment: utils.Ptr("comment"),
		Records: &[]dns.RecordPayload{
			{Content: utils.Ptr("1.1.1.1")},
		},
		Ttl: utils.Ptr(int64(3600)),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description      string
		argValues        []string
		flagValues       map[string]string
		recordFlagValues []string
		isValid          bool
		expectedModel    *inputModel
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
				zoneIdFlag:    testZoneId,
			},
			isValid: false,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				ZoneId:      testZoneId,
				RecordSetId: testRecordSetId,
			},
		},
		{
			description: "zero values",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
				zoneIdFlag:    testZoneId,
				commentFlag:   "",
				nameFlag:      "",
				recordFlag:    "1.1.1.1",
				ttlFlag:       "0",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				ZoneId:      testZoneId,
				RecordSetId: testRecordSetId,
				Name:        utils.Ptr(""),
				Comment:     utils.Ptr(""),
				Records:     &[]string{"1.1.1.1"},
				TTL:         utils.Ptr(int64(0)),
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
			description: "zone id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, zoneIdFlag)
			}),
			isValid: false,
		},
		{
			description: "zone id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[zoneIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "zone id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[zoneIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "record set id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "record set id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description:      "repeated primary flags",
			argValues:        fixtureArgValues(),
			flagValues:       fixtureFlagValues(),
			recordFlagValues: []string{"1.2.3.4", "5.6.7.8"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Records = utils.Ptr(append(*model.Records, "1.2.3.4", "5.6.7.8"))
			}),
		},
		{
			description:      "repeated primary flags with list value",
			argValues:        fixtureArgValues(),
			flagValues:       fixtureFlagValues(),
			recordFlagValues: []string{"1.2.3.4,5.6.7.8"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Records = utils.Ptr(append(*model.Records, "1.2.3.4", "5.6.7.8"))
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

			for _, value := range tt.recordFlagValues {
				err := cmd.Flags().Set(recordFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", recordFlag, value, err)
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
		expectedRequest dns.ApiPartialUpdateRecordSetRequest
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
				},
				ZoneId:      testZoneId,
				RecordSetId: testRecordSetId,
			},
			expectedRequest: testClient.PartialUpdateRecordSet(testCtx, testProjectId, testZoneId, testRecordSetId).
				PartialUpdateRecordSetPayload(dns.PartialUpdateRecordSetPayload{}),
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
