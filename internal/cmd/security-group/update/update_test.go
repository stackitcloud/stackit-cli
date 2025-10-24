package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var (
	testCtx         = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient      = &iaas.APIClient{}
	testProjectId   = uuid.NewString()
	testGroupId     = []string{uuid.NewString()}
	testName        = "new-security-group"
	testDescription = "a test description"
	testLabels      = map[string]string{
		"fooKey": "fooValue",
		"barKey": "barValue",
		"bazKey": "bazValue",
	}
)

func toStringAnyMapPtr(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	result := map[string]any{}
	for k, v := range m {
		result[k] = v
	}
	return result
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		descriptionArg: testDescription,
		labelsArg:      "fooKey=fooValue,barKey=barValue,bazKey=bazValue",
		nameArg:        testName,
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
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		Labels:          &testLabels,
		Description:     &testDescription,
		Name:            &testName,
		SecurityGroupId: testGroupId[0],
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiUpdateSecurityGroupRequest)) iaas.ApiUpdateSecurityGroupRequest {
	request := testClient.UpdateSecurityGroup(testCtx, testProjectId, testRegion, testGroupId[0])
	request = request.UpdateSecurityGroupPayload(iaas.UpdateSecurityGroupPayload{
		Description: &testDescription,
		Labels:      utils.Ptr(toStringAnyMapPtr(testLabels)),
		Name:        &testName,
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		args          []string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			args:          testGroupId,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values but valid group id",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
			},
			args:    testGroupId,
			isValid: false,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
				model.Name = nil
				model.Description = nil
			}),
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			args:    testGroupId,
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			args:    testGroupId,
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			args:    testGroupId,
			isValid: false,
		},
		{
			description: "no name passed",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameArg)
			}),
			args: testGroupId,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = nil
			}),
			isValid: true,
		},
		{
			description: "no description passed",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, descriptionArg)
			}),
			args: testGroupId,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
			isValid: true,
		},
		{
			description: "no labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelsArg)
			}),
			args: testGroupId,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
			isValid: true,
		},
		{
			description: "single label",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelsArg] = "foo=bar"
			}),
			args:    testGroupId,
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{
					"foo": "bar",
				}
			}),
		},
		{
			description: "no group id passed",
			flagValues:  fixtureFlagValues(),
			args:        nil,
			isValid:     false,
		},
		{
			description: "invalid group id passed",
			flagValues:  fixtureFlagValues(),
			args:        []string{"foobar"},
			isValid:     false,
		},
		{
			description: "multiple group ids passed",
			flagValues:  fixtureFlagValues(),
			args:        []string{uuid.NewString(), uuid.NewString()},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			if err := globalflags.Configure(cmd.Flags()); err != nil {
				t.Errorf("cannot configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				if err := cmd.Flags().Set(flag, value); err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			if err := cmd.ValidateRequiredFlags(); err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			if err := cmd.ValidateArgs(tt.args); err != nil {
				if !tt.isValid {
					return
				}
			}

			model, err := parseInput(p, cmd, tt.args)
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
		expectedRequest iaas.ApiUpdateSecurityGroupRequest
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
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateSecurityGroupRequest) {
				*request = (*request).UpdateSecurityGroupPayload(iaas.UpdateSecurityGroupPayload{
					Description: &testDescription,
					Labels:      nil,
					Name:        &testName,
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
