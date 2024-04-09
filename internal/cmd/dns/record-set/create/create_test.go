package create

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

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		zoneIdFlag:    testZoneId,
		commentFlag:   "comment",
		nameFlag:      "example.com",
		recordFlag:    "1.1.1.1",
		ttlFlag:       "3600",
		typeFlag:      "SOA", // Non-default value
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
		ZoneId:  testZoneId,
		Name:    utils.Ptr("example.com"),
		Comment: utils.Ptr("comment"),
		Records: []string{"1.1.1.1"},
		TTL:     utils.Ptr(int64(3600)),
		Type:    "SOA",
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiCreateRecordSetRequest)) dns.ApiCreateRecordSetRequest {
	request := testClient.CreateRecordSet(testCtx, testProjectId, testZoneId)
	request = request.CreateRecordSetPayload(dns.CreateRecordSetPayload{
		Name:    utils.Ptr("example.com"),
		Comment: utils.Ptr("comment"),
		Records: &[]dns.RecordPayload{
			{Content: utils.Ptr("1.1.1.1")},
		},
		Ttl:  utils.Ptr(int64(3600)),
		Type: utils.Ptr("SOA"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description      string
		flagValues       map[string]string
		recordFlagValues []string
		isValid          bool
		expectedModel    *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
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
				zoneIdFlag:    testZoneId,
				nameFlag:      "example.com",
				recordFlag:    "1.1.1.1",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				ZoneId:  testZoneId,
				Name:    utils.Ptr("example.com"),
				Records: []string{"1.1.1.1"},
				Type:    defaultType,
			},
		},
		{
			description: "zero values",
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
					Verbosity: globalflags.VerbosityDefault,
				},
				ZoneId:  testZoneId,
				Name:    utils.Ptr(""),
				Comment: utils.Ptr(""),
				Records: []string{"1.1.1.1"},
				TTL:     utils.Ptr(int64(0)),
				Type:    defaultType,
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
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "records missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, recordFlag)
			}),
			isValid: false,
		},
		{
			description: "type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, typeFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Type = defaultType
			}),
		},
		{
			description: "type invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "type invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag] = "a"
			}),
			isValid: false,
		},
		{
			description:      "repeated primary flags",
			flagValues:       fixtureFlagValues(),
			recordFlagValues: []string{"1.2.3.4", "5.6.7.8"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Records = append(model.Records, "1.2.3.4", "5.6.7.8")
			}),
		},
		{
			description:      "repeated primary flags with list value",
			flagValues:       fixtureFlagValues(),
			recordFlagValues: []string{"1.2.3.4,5.6.7.8"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Records = append(model.Records, "1.2.3.4", "5.6.7.8")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd(nil)
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

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(cmd)
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
		expectedRequest dns.ApiCreateRecordSetRequest
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
				ZoneId:  testZoneId,
				Name:    utils.Ptr("example.com"),
				Records: []string{"1.1.1.1"},
				Type:    defaultType,
			},
			expectedRequest: testClient.CreateRecordSet(testCtx, testProjectId, testZoneId).
				CreateRecordSetPayload(dns.CreateRecordSetPayload{
					Name: utils.Ptr("example.com"),
					Records: &[]dns.RecordPayload{
						{Content: utils.Ptr("1.1.1.1")},
					},
					Type: utils.Ptr(defaultType),
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
