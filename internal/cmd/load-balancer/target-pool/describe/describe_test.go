package describe

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &loadbalancer.APIClient{}
	testProjectId = uuid.NewString()
)

const (
	testLoadBalancerName = "my-load-balancer"
	testTargetPoolName   = "target-pool-1"
	testTargetName       = "my-target"
	testIp               = "1.2.3.4"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testTargetPoolName,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:        testProjectId,
		loadBalancerNameFlag: testLoadBalancerName,
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
		LoadBalancerName: testLoadBalancerName,
		TargetPoolName:   testTargetPoolName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureTargets() *[]loadbalancer.Target {
	return &[]loadbalancer.Target{
		{
			DisplayName: utils.Ptr("target-1"),
			Ip:          utils.Ptr("1.2.3.4"),
		},
		{
			DisplayName: utils.Ptr("target-2"),
			Ip:          utils.Ptr("4.3.2.1"),
		},
	}
}

func fixtureLoadBalancer(mods ...func(*loadbalancer.LoadBalancer)) *loadbalancer.LoadBalancer {
	lb := loadbalancer.LoadBalancer{
		Name: utils.Ptr(testLoadBalancerName),
		TargetPools: &[]loadbalancer.TargetPool{
			{
				Name:    utils.Ptr(testTargetPoolName),
				Targets: fixtureTargets(),
				ActiveHealthCheck: &loadbalancer.ActiveHealthCheck{
					UnhealthyThreshold: utils.Ptr(int64(3)),
				},
				SessionPersistence: &loadbalancer.SessionPersistence{
					UseSourceIpAddress: utils.Ptr(true),
				},
				TargetPort: utils.Ptr(int64(80)),
			},
			{
				Name: utils.Ptr("target-pool-2"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("6.7.8.9"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("9.8.7.6"),
					},
				},
			},
		},
	}

	for _, mod := range mods {
		mod(&lb)
	}
	return &lb
}

func fixtureRequest(mods ...func(request *loadbalancer.ApiGetLoadBalancerRequest)) loadbalancer.ApiGetLoadBalancerRequest {
	request := testClient.GetLoadBalancer(testCtx, testProjectId, testLoadBalancerName)
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "target pool name empty",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "load balancer name missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, loadBalancerNameFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)

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

			err = cmd.ValidateArgs(tt.argValues)
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

			model, err := parseInput(p, cmd, tt.argValues)
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
		isValid         bool
		expectedRequest loadbalancer.ApiGetLoadBalancerRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
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
