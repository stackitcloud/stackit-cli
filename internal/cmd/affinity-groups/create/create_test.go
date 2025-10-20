package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &iaas.APIClient{}
	testProjectId = uuid.NewString()
)

const (
	testName   = "test-name"
	testPolicy = "test-policy"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,

		nameFlag:   testName,
		policyFlag: testPolicy,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
		},
		Name:   testName,
		Policy: testPolicy,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateAffinityGroupRequest)) iaas.ApiCreateAffinityGroupRequest {
	request := testClient.CreateAffinityGroup(testCtx, testProjectId)
	request = request.CreateAffinityGroupPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateAffinityGroupPayload)) iaas.CreateAffinityGroupPayload {
	payload := iaas.CreateAffinityGroupPayload{
		Name:   utils.Ptr(testName),
		Policy: utils.Ptr(testPolicy),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			description: "without name flag",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, "name")
				},
			),
			isValid: false,
		},
		{
			description: "without policy flag",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, "policy")
				},
			),
			isValid: false,
		},
		{
			description: "without name and policy flag",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, "policy")
					delete(flagValues, "name")
				},
			),
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
		model           inputModel
		expectedRequest iaas.ApiCreateAffinityGroupRequest
	}{
		{
			description:     "base",
			model:           *fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx))
			if diff != "" {
				t.Fatalf("Request does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description string
		model       inputModel
		response    iaas.AffinityGroup
		isValid     bool
	}{
		{
			description: "empty",
			model:       inputModel{},
			response:    iaas.AffinityGroup{},
			isValid:     true,
		},
		{
			description: "base",
			model:       *fixtureInputModel(),
			response: iaas.AffinityGroup{
				Id:      utils.Ptr(testProjectId),
				Members: utils.Ptr([]string{uuid.NewString(), uuid.NewString()}),
				Name:    utils.Ptr("test-project"),
				Policy:  utils.Ptr("hard-affinity"),
			},
			isValid: true,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(p, tt.model, tt.response)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error output result: %v", err)
				return
			}
			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
				return
			}
		})
	}
}
