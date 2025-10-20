package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &resourcemanager.APIClient{}
var testProjectId = uuid.NewString()
var testParentId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		parentIdFlag:  testParentId,
		nameFlag:      nameFlag,
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
		ParentId: utils.Ptr(testParentId),
		Name:     utils.Ptr(nameFlag),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *resourcemanager.ApiPartialUpdateProjectRequest)) resourcemanager.ApiPartialUpdateProjectRequest {
	request := testClient.PartialUpdateProject(testCtx, testProjectId)
	request = request.PartialUpdateProjectPayload(resourcemanager.PartialUpdateProjectPayload{
		ContainerParentId: utils.Ptr(testParentId),
		Name:              utils.Ptr(nameFlag),
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
		labelValues   []string
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
			description: "required flags only (no values to update)",
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
			},
			isValid: false,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
			},
		},
		{
			description: "valid_labels",
			flagValues:  fixtureFlagValues(),
			labelValues: []string{"key=value", "foo=bar"},
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.Labels = &map[string]string{
						"key": "value",
						"foo": "bar",
					}
				}),
			isValid: true,
		},
		{
			description: "valid_labels_2",
			flagValues:  fixtureFlagValues(),
			labelValues: []string{"key=value,foo=bar"},
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.Labels = &map[string]string{
						"key": "value",
						"foo": "bar",
					}
				}),
			isValid: true,
		},
		{
			description: "invalid_labels",
			flagValues:  fixtureFlagValues(),
			labelValues: []string{"key"},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				labelFlag: tt.labelValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest resourcemanager.ApiPartialUpdateProjectRequest
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
			},
			expectedRequest: testClient.PartialUpdateProject(testCtx, testProjectId).
				PartialUpdateProjectPayload(resourcemanager.PartialUpdateProjectPayload{}),
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
