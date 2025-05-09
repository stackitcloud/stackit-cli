package addtarget

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &loadbalancer.APIClient{}
	testProjectId = uuid.NewString()
)

const (
	testRegion         = "eu02"
	testLBName         = "my-load-balancer"
	testTargetPoolName = "target-pool-1"
	testTargetName     = "my-target"
	testIP             = "1.1.1.1"
)

type loadBalancerClientMocked struct {
	getCredentialsFails  bool
	getCredentialsResp   *loadbalancer.GetCredentialsResponse
	getLoadBalancerFails bool
	getLoadBalancerResp  *loadbalancer.LoadBalancer
}

func (m *loadBalancerClientMocked) GetCredentialsExecute(_ context.Context, _, _, _ string) (*loadbalancer.GetCredentialsResponse, error) {
	if m.getCredentialsFails {
		return nil, fmt.Errorf("could not get credentials")
	}
	return m.getCredentialsResp, nil
}

func (m *loadBalancerClientMocked) GetLoadBalancerExecute(_ context.Context, _, _, _ string) (*loadbalancer.LoadBalancer, error) {
	if m.getLoadBalancerFails {
		return nil, fmt.Errorf("could not get load balancer")
	}
	return m.getLoadBalancerResp, nil
}

func (m *loadBalancerClientMocked) UpdateTargetPool(ctx context.Context, projectId, region, loadBalancerName, targetPoolName string) loadbalancer.ApiUpdateTargetPoolRequest {
	return testClient.UpdateTargetPool(ctx, projectId, region, loadBalancerName, targetPoolName)
}

func (m *loadBalancerClientMocked) ListLoadBalancersExecute(_ context.Context, _, _ string) (*loadbalancer.ListLoadBalancersResponse, error) {
	return nil, nil
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testIP,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		lbNameFlag:                testLBName,
		targetNameFlag:            testTargetName,
		targetPoolNameFlag:        testTargetPoolName,
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
		TargetPoolName: testTargetPoolName,
		LBName:         testLBName,
		TargetName:     testTargetName,
		IP:             testIP,
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
		Name: utils.Ptr(testLBName),
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

func fixturePayload(mods ...func(payload *loadbalancer.UpdateTargetPoolPayload)) *loadbalancer.UpdateTargetPoolPayload {
	payload := &loadbalancer.UpdateTargetPoolPayload{
		Name: utils.Ptr("target-pool-1"),
		ActiveHealthCheck: &loadbalancer.ActiveHealthCheck{
			UnhealthyThreshold: utils.Ptr(int64(3)),
		},
		SessionPersistence: &loadbalancer.SessionPersistence{
			UseSourceIpAddress: utils.Ptr(true),
		},
		TargetPort: utils.Ptr(int64(80)),
		Targets:    fixtureTargets(),
	}

	for _, mod := range mods {
		mod(payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *loadbalancer.ApiUpdateTargetPoolRequest)) loadbalancer.ApiUpdateTargetPoolRequest {
	request := testClient.UpdateTargetPool(testCtx, testProjectId, testRegion, testLBName, testTargetPoolName)
	request = request.UpdateTargetPoolPayload(*fixturePayload())
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
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "ip missing",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "load balancer name missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, lbNameFlag)
			}),
			isValid: false,
		},
		{
			description: "target name missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, targetNameFlag)
			}),
			isValid: false,
		},
		{
			description: "target pool name missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, targetPoolNameFlag)
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
		description          string
		model                *inputModel
		isValid              bool
		getLoadBalancerFails bool
		getLoadBalancerResp  *loadbalancer.LoadBalancer
		expectedRequest      loadbalancer.ApiUpdateTargetPoolRequest
	}{
		{
			description:         "base",
			model:               fixtureInputModel(),
			getLoadBalancerResp: fixtureLoadBalancer(),
			isValid:             true,
			expectedRequest: fixtureRequest(func(request *loadbalancer.ApiUpdateTargetPoolRequest) {
				payload := fixturePayload(func(payload *loadbalancer.UpdateTargetPoolPayload) {
					payload.Targets = &[]loadbalancer.Target{
						(*fixtureTargets())[0],
						(*fixtureTargets())[1],
						{
							DisplayName: utils.Ptr(testTargetName),
							Ip:          utils.Ptr(testIP),
						},
					}
				})
				*request = request.UpdateTargetPoolPayload(*payload)
			}),
		},
		{
			description: "empty targets",
			model:       fixtureInputModel(),
			getLoadBalancerResp: fixtureLoadBalancer(func(lb *loadbalancer.LoadBalancer) {
				(*lb.TargetPools)[0].Targets = &[]loadbalancer.Target{}
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *loadbalancer.ApiUpdateTargetPoolRequest) {
				payload := fixturePayload(func(payload *loadbalancer.UpdateTargetPoolPayload) {
					payload.Targets = &[]loadbalancer.Target{
						{
							DisplayName: utils.Ptr(testTargetName),
							Ip:          utils.Ptr(testIP),
						},
					}
				})
				*request = request.UpdateTargetPoolPayload(*payload)
			}),
		},
		{
			description: "nil targets",
			model:       fixtureInputModel(),
			getLoadBalancerResp: fixtureLoadBalancer(func(lb *loadbalancer.LoadBalancer) {
				(*lb.TargetPools)[0].Targets = nil
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *loadbalancer.ApiUpdateTargetPoolRequest) {
				payload := fixturePayload(func(payload *loadbalancer.UpdateTargetPoolPayload) {
					payload.Targets = &[]loadbalancer.Target{
						{
							DisplayName: utils.Ptr(testTargetName),
							Ip:          utils.Ptr(testIP),
						},
					}
				})
				*request = request.UpdateTargetPoolPayload(*payload)
			}),
		},
		{
			description:          "get load balancer fails",
			model:                fixtureInputModel(),
			getLoadBalancerFails: true,
			isValid:              false,
		},
		{
			description: "target pool not found",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.TargetPoolName = "not-existent"
				}),
			getLoadBalancerResp: fixtureLoadBalancer(),
			isValid:             false,
		},
		{
			description: "nil target pool",
			model:       fixtureInputModel(),
			getLoadBalancerResp: fixtureLoadBalancer(func(lb *loadbalancer.LoadBalancer) {
				*lb.TargetPools = nil
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &loadBalancerClientMocked{
				getLoadBalancerFails: tt.getLoadBalancerFails,
				getLoadBalancerResp:  tt.getLoadBalancerResp,
			}
			request, err := buildRequest(testCtx, tt.model, client)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
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
