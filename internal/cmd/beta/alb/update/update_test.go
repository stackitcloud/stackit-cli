package update

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

//go:embed testdata/testconfig.json
var testConfiguration []byte

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var (
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient       = &alb.APIClient{}
	testProjectId    = uuid.NewString()
	testRegion       = "eu01"
	testLoadBalancer = "my-load-balancer"
	testConfig       = "testdata/testconfig.json"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		configurationFlag:      testConfig,
		globalflags.RegionFlag: testRegion,
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
			Region:    testRegion,
		},
		Configuration: utils.Ptr(testConfig),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}
func fixturePayload(mods ...func(payload *alb.UpdateLoadBalancerPayload)) (payload alb.UpdateLoadBalancerPayload) {
	if err := json.Unmarshal(testConfiguration, &payload); err != nil {
		log.Panicf("cannot deserialize test configuration: %v", err)
	}
	for _, f := range mods {
		f(&payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *alb.ApiUpdateLoadBalancerRequest)) alb.ApiUpdateLoadBalancerRequest {
	request := testClient.UpdateLoadBalancer(testCtx, testProjectId, testRegion, testLoadBalancer)

	request = request.UpdateLoadBalancerPayload(fixturePayload())
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
			description: "required fields only",
			flagValues: map[string]string{
				projectIdFlag:     testProjectId,
				configurationFlag: testConfig,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				Configuration: &testConfig,
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
		expectedRequest alb.ApiUpdateLoadBalancerRequest
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
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				Configuration: &testConfig,
			},
			expectedRequest: testClient.
				UpdateLoadBalancer(testCtx, testProjectId, testRegion, testLoadBalancer).
				UpdateLoadBalancerPayload(fixturePayload()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				t.Fatalf("cannot create request: %v", err)
			}

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
		model        *inputModel
		projectLabel string
		resp         *alb.LoadBalancer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "empty response as argument",
			args: args{
				model: fixtureInputModel(),
				resp:  &alb.LoadBalancer{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
