package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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
	testName        = "new-security-group"
	testDescription = "a test description"
	testLabels      = map[string]string{
		"fooKey": "fooValue",
		"barKey": "barValue",
		"bazKey": "bazValue",
	}
	testStateful = true
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		descriptionFlag: testDescription,
		labelsFlag:      "fooKey=fooValue,barKey=barValue,bazKey=bazValue",
		statefulFlag:    "true",
		nameFlag:        testName,
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
		Labels:      &testLabels,
		Description: &testDescription,
		Name:        &testName,
		Stateful:    &testStateful,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

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
func fixtureRequest(mods ...func(request *iaas.ApiCreateSecurityGroupRequest)) iaas.ApiCreateSecurityGroupRequest {
	request := testClient.CreateSecurityGroup(testCtx, testProjectId, testRegion)

	request = request.CreateSecurityGroupPayload(iaas.CreateSecurityGroupPayload{
		Description: &testDescription,
		Labels:      utils.Ptr(toStringAnyMapPtr(testLabels)),
		Name:        &testName,
		Rules:       nil,
		Stateful:    &testStateful,
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
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
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
			description: "no labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelsFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
		},
		{
			description: "single label",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelsFlag] = "foo=bar"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{
					"foo": "bar",
				}
			}),
		},
		{
			description: "stateless security group",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[statefulFlag] = "false"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Stateful = utils.Ptr(false)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
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
				*request = (*request).CreateSecurityGroupPayload(iaas.CreateSecurityGroupPayload{
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
				model.Stateful = utils.Ptr(false)
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateSecurityGroupRequest) {
				*request = (*request).CreateSecurityGroupPayload(iaas.CreateSecurityGroupPayload{
					Description: &testDescription,
					Labels:      utils.Ptr(toStringAnyMapPtr(testLabels)),
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

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		resp         iaas.SecurityGroup
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.name, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
