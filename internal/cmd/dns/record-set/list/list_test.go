package list

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{}
var testProjectId = uuid.NewString()
var testZoneId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:   testProjectId,
		zoneIdFlag:      testZoneId,
		nameLikeFlag:    "some-pattern",
		activeFlag:      "true",
		orderByNameFlag: "asc",
		limitFlag:       "10",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagModel(mods ...func(model *flagModel)) *flagModel {
	model := &flagModel{
		ProjectId:   testProjectId,
		ZoneId:      testZoneId,
		NameLike:    utils.Ptr("some-pattern"),
		Active:      utils.Ptr(true),
		OrderByName: utils.Ptr("asc"),
		Limit:       utils.Ptr(int64(10)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiGetRecordSetsRequest)) dns.ApiGetRecordSetsRequest {
	request := testClient.GetRecordSets(testCtx, testProjectId, testZoneId)
	request = request.NameLike("some-pattern")
	request = request.ActiveEq(true)
	request = request.OrderByName("ASC")
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		isValid       bool
		expectedModel *flagModel
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
				zoneIdFlag:    testZoneId,
			},
			isValid: true,
			expectedModel: &flagModel{
				ProjectId: testProjectId,
				ZoneId:    testZoneId,
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
			description: "name like empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameLikeFlag] = ""
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.NameLike = utils.Ptr("")
			}),
		},
		{
			description: "is active = false",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = "false"
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.Active = utils.Ptr(false)
			}),
		},
		{
			description: "is active invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "is active invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "order by name desc",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[orderByNameFlag] = "desc"
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.OrderByName = utils.Ptr("desc")
			}),
		},
		{
			description: "order by name invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[orderByNameFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "order by name invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[orderByNameFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "limit invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "limit invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Flag defined in root command
			err := testutils.ConfigureBindUUIDFlag(cmd, projectIdFlag, config.ProjectIdKey)
			if err != nil {
				t.Fatalf("configure global flag --%s: %v", projectIdFlag, err)
			}

			configureFlags(cmd)

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
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
		expectedRequest dns.ApiGetRecordSetsRequest
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
			expectedRequest: testClient.GetRecordSets(testCtx, testProjectId, testZoneId),
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
