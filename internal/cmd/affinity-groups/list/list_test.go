package list

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
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
	testLimit = 10
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
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
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiListAffinityGroupsRequest)) iaas.ApiListAffinityGroupsRequest {
	request := testClient.ListAffinityGroups(testCtx, testProjectId)
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
			description: "without flags",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "with limit flag",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["limit"] = strconv.Itoa(testLimit)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(testLimit))
			}),
		},
		{
			description: "with limit flag == 0",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["limit"] = strconv.Itoa(0)
			}),
			isValid: false,
		},
		{
			description: "with limit flag < 0",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["limit"] = strconv.Itoa(-1)
			}),
			isValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			if err := globalflags.Configure(cmd.Flags()); err != nil {
				t.Fatalf("configure global flags: %v", err)
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
		model           inputModel
		expectedRequest iaas.ApiListAffinityGroupsRequest
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
		response    []iaas.AffinityGroup
		isValid     bool
	}{
		{
			description: "empty",
			model:       inputModel{},
			response:    []iaas.AffinityGroup{},
			isValid:     true,
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
