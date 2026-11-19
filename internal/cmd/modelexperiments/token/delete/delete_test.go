package delete

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testContextKey struct{}

const testRegion = "eu01"

var (
	testContext    = context.WithValue(context.Background(), testContextKey{}, "token-delete")
	testClient     = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
	testProjectID  = uuid.NewString()
	testInstanceID = uuid.NewString()
	testTokenID    = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,
		instanceIDFlag:            testInstanceID,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{testTokenID}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectID,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		TokenID:    testTokenID,
		InstanceID: testInstanceID,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest() modelexperiments.ApiDeleteInstanceTokenRequest {
	return testClient.DefaultAPI.DeleteInstanceToken(testContext, testProjectID, testRegion, testTokenID, testInstanceID)
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
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIDFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id empty",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "instance id invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "token id empty",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "token id invalid",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	request := buildRequest(testContext, fixtureInputModel(), testClient)
	expectedRequest := fixtureRequest()

	diff := cmp.Diff(request, expectedRequest,
		cmp.AllowUnexported(expectedRequest),
		cmpopts.EquateComparable(testContext),
	)
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}
