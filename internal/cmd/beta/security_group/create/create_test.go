package create

import (
	"context"
	"strings"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var (
	testCtx         = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient      = &iaas.APIClient{}
	testProjectId   = uuid.NewString()
	testName        = "new-security-group"
	testDescription = "a test description"
	testLabels      = map[string]any{
		"fooKey": "fooValue",
		"barKey": "barValue",
		"bazKey": "bazValue",
	}
	testStateful = true
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		"description": testDescription,
		"labels":      "fooKey=fooValue,barKey=barValue,bazKey=bazValue",
		"stateful":    "true",
		"name":        testName,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectId, Verbosity: globalflags.VerbosityDefault},
		Labels:          testLabels,
		Description:     testDescription,
		Name:            testName,
		Stateful:        testStateful,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateSecurityGroupRequest)) iaas.ApiCreateSecurityGroupRequest {
	request := testClient.CreateSecurityGroup(testCtx, testProjectId)
	request = request.CreateSecurityGroupPayload(iaas.CreateSecurityGroupPayload{
		Description: utils.Ptr(testDescription),
		Labels:      &testLabels,
		Name:        utils.Ptr(testName),
		Rules:       nil,
		Stateful:    utils.Ptr(testStateful),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
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
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, "name")
			}),
			isValid: false,
		},
		{
			description: "name too long",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["name"] = strings.Repeat("toolong", 1000)
			}),
			isValid: false,
		},
		{
			description: "description too long",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["description"] = strings.Repeat("toolong", 1000)
			}),
			isValid: false,
		},
		{
			description: "no labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, "labels")
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = map[string]any{}
			}),
		},
		{
			description: "single label",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["labels"] = "foo=bar"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = map[string]any{
					"foo": "bar",
				}
			}),
		},
		{
			description: "malformed labels 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["labels"] = "foo=bar=baz"
			}),
			isValid: false,
		},
		{
			description: "malformed labels 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["labels"] = "foobarbaz"
			}),
			isValid: false,
		},
		{
			description: "stateless security group",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["stateful"] = "false"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Stateful = false
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

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			p := print.NewPrinter()
			model, err := parseInput(p, cmd)
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
		expectedRequest iaas.ApiCreateSecurityGroupRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no labels",
			model: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateSecurityGroupRequest) {
				*request = request.CreateSecurityGroupPayload(iaas.CreateSecurityGroupPayload{
					Description: &testDescription,
					Labels:      nil,
					Name:        &testName,
					Stateful:    &testStateful,
				})
			}),
		},
		{
			description: "stateless security group",
			model: fixtureInputModel(func(model *inputModel) {
				model.Stateful = false
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateSecurityGroupRequest) {
				*request = request.CreateSecurityGroupPayload(iaas.CreateSecurityGroupPayload{
					Description: &testDescription,
					Labels:      &testLabels,
					Name:        &testName,
					Stateful:    utils.Ptr(false),
				})
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
