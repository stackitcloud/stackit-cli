package update

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testContextKey struct{}

const (
	testRegion = "eu01"
	testName   = "updated"
)

var (
	testContext    = context.WithValue(context.Background(), testContextKey{}, "token-update")
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

func fixtureRequest(payload modelexperiments.PartialUpdateInstanceTokenPayload) modelexperiments.ApiPartialUpdateInstanceTokenRequest {
	return testClient.DefaultAPI.PartialUpdateInstanceToken(testContext, testProjectID, testRegion, testTokenID, testInstanceID).
		PartialUpdateInstanceTokenPayload(payload)
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
			description: "name",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameFlag] = testName
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = utils.Ptr(testName)
			}),
		},
		{
			description: "description",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[descriptionFlag] = "updated token"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("updated token")
			}),
		},
		{
			description: "labels",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelsFlag] = "env=prod,team=ml"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]*string{
					"env":  utils.Ptr("prod"),
					"team": utils.Ptr("ml"),
				}
			}),
		},
		{
			description: "no update flags",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameFlag] = testName
			}),
			isValid: false,
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
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
		{
			description: "instance id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIDFlag)
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
		{
			description: "instance id empty",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = ""
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
		{
			description: "instance id invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = "invalid-uuid"
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
		{
			description: "token id empty",
			argValues:   []string{""},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
		{
			description: "token id invalid",
			argValues:   []string{"invalid-uuid"},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameFlag] = testName
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	labels := map[string]*string{
		"env":  utils.Ptr("prod"),
		"team": utils.Ptr("ml"),
	}
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest modelexperiments.ApiPartialUpdateInstanceTokenRequest
	}{
		{
			description: "name",
			model: fixtureInputModel(func(model *inputModel) {
				model.Name = utils.Ptr(testName)
			}),
			expectedRequest: fixtureRequest(modelexperiments.PartialUpdateInstanceTokenPayload{
				Name: utils.Ptr(testName),
			}),
		},
		{
			description: "description",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("updated token")
			}),
			expectedRequest: fixtureRequest(modelexperiments.PartialUpdateInstanceTokenPayload{
				Description: utils.Ptr("updated token"),
			}),
		},
		{
			description: "labels",
			model: fixtureInputModel(func(model *inputModel) {
				model.Labels = &labels
			}),
			expectedRequest: fixtureRequest(modelexperiments.PartialUpdateInstanceTokenPayload{
				Labels: &labels,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := buildRequest(testContext, tt.model, testClient)
			diff := cmp.Diff(got, tt.expectedRequest, cmp.AllowUnexported(tt.expectedRequest), cmpopts.EquateComparable(testContext))
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name    string
		resp    *modelexperiments.PartialUpdateInstanceTokenResponse
		wantErr bool
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name:    "empty token",
			resp:    &modelexperiments.PartialUpdateInstanceTokenResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()
			if err := outputResult(params.Printer, fixtureInputModel(), tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
