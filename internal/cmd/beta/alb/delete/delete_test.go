package delete

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

type testCtxKey struct{}

var (
	testCtx              = context.WithValue(context.Background(), testCtxKey{}, "test")
	testProjectId        = uuid.NewString()
	testRegion           = "eu01"
	testClient           = &alb.APIClient{}
	testLoadBalancerName = "my-test-loadbalancer"
)

func fixtureArgValues(mods ...func(argVales []string)) []string {
	argVales := []string{
		testLoadBalancerName,
	}
	for _, m := range mods {
		m(argVales)
	}
	return argVales
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
	}
	for _, m := range mods {
		m(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
			Region:    testRegion,
		},
		Name: testLoadBalancerName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *alb.ApiDeleteLoadBalancerRequest)) alb.ApiDeleteLoadBalancerRequest {
	request := testClient.DeleteLoadBalancer(testCtx, testProjectId, testRegion, testLoadBalancerName)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argsValues    []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argsValues:    fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argsValues:  []string{},
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
			},
			isValid: false,
		},
		{
			description: "no arg values",
			argsValues:  []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argsValues:  fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
			},
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err = cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateArgs(tt.argsValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argsValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
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
		description    string
		model          *inputModel
		expectedResult alb.ApiDeleteLoadBalancerRequest
	}{
		{
			description:    "base",
			model:          fixtureInputModel(),
			expectedResult: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedResult,
				cmp.AllowUnexported(tt.expectedResult),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}
