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
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	testRegion           = "eu02"
	testLoadBalancerName = "loadBalancer"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &loadbalancer.APIClient{}
var testProjectId = uuid.NewString()

var testPayload = loadbalancer.UpdateLoadBalancerPayload{
	ExternalAddress: utils.Ptr(""),

	Listeners: &[]loadbalancer.Listener{
		{
			DisplayName: utils.Ptr(""),
			Port:        utils.Ptr(int64(0)),
			Protocol:    loadbalancer.ListenerProtocol("").Ptr(),
			ServerNameIndicators: &[]loadbalancer.ServerNameIndicator{
				{
					Name: utils.Ptr(""),
				},
			},
			TargetPool: utils.Ptr(""),
			Tcp: &loadbalancer.OptionsTCP{
				IdleTimeout: utils.Ptr(""),
			},
			Udp: &loadbalancer.OptionsUDP{
				IdleTimeout: utils.Ptr(""),
			},
		},
	},
	Name: utils.Ptr(""),
	Networks: &[]loadbalancer.Network{
		{
			NetworkId: utils.Ptr(""),
			Role:      loadbalancer.NetworkRole("").Ptr(),
		},
	},
	Options: &loadbalancer.LoadBalancerOptions{
		AccessControl: &loadbalancer.LoadbalancerOptionAccessControl{
			AllowedSourceRanges: &[]string{
				"",
			},
		},
		EphemeralAddress: utils.Ptr(false),
		Observability: &loadbalancer.LoadbalancerOptionObservability{
			Logs: &loadbalancer.LoadbalancerOptionLogs{
				CredentialsRef: utils.Ptr(""),
				PushUrl:        utils.Ptr(""),
			},
			Metrics: &loadbalancer.LoadbalancerOptionMetrics{
				CredentialsRef: utils.Ptr(""),
				PushUrl:        utils.Ptr(""),
			},
		},
		PrivateNetworkOnly: utils.Ptr(false),
	},
	TargetPools: &[]loadbalancer.TargetPool{
		{
			ActiveHealthCheck: &loadbalancer.ActiveHealthCheck{
				HealthyThreshold:   utils.Ptr(int64(0)),
				Interval:           utils.Ptr(""),
				IntervalJitter:     utils.Ptr(""),
				Timeout:            utils.Ptr(""),
				UnhealthyThreshold: utils.Ptr(int64(0)),
			},
			Name: utils.Ptr(""),
			SessionPersistence: &loadbalancer.SessionPersistence{
				UseSourceIpAddress: utils.Ptr(false),
			},
			TargetPort: utils.Ptr(int64(0)),
			Targets: &[]loadbalancer.Target{
				{
					DisplayName: utils.Ptr(""),
					Ip:          utils.Ptr(""),
				},
			},
		},
	},
	Version: utils.Ptr(""),
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testLoadBalancerName,
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
		payloadFlag: `
{
  "externalAddress": "",
  "listeners": [
    {
      "displayName": "",
      "port": 0,
      "protocol": "",
      "serverNameIndicators": [
        {
          "name": ""
        }
      ],
      "targetPool": "",
      "tcp": {
        "idleTimeout": ""
      },
      "udp": {
        "idleTimeout": ""
      }
    }
  ],
  "name": "",
  "networks": [
    {
      "networkId": "",
      "role": ""
    }
  ],
  "options": {
    "accessControl": {
      "allowedSourceRanges": [
        ""
      ]
    },
    "ephemeralAddress": false,
    "observability": {
      "logs": {
        "credentialsRef": "",
        "pushUrl": ""
      },
      "metrics": {
        "credentialsRef": "",
        "pushUrl": ""
      }
    },
    "privateNetworkOnly": false
  },
  "targetPools": [
    {
      "activeHealthCheck": {
        "healthyThreshold": 0,
        "interval": "",
        "intervalJitter": "",
        "timeout": "",
        "unhealthyThreshold": 0
      },
      "name": "",
      "sessionPersistence": {
        "useSourceIpAddress": false
      },
      "targetPort": 0,
      "targets": [
        {
          "displayName": "",
          "ip": ""
        }
      ]
    }
  ],
  "version": ""
}
`,
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
		LoadBalancerName: testLoadBalancerName,
		Payload:          testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *loadbalancer.ApiUpdateLoadBalancerRequest)) loadbalancer.ApiUpdateLoadBalancerRequest {
	request := testClient.UpdateLoadBalancer(testCtx, testProjectId, testRegion, testLoadBalancerName)
	request = request.UpdateLoadBalancerPayload(testPayload)
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
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
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
			description: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid: false,
		},
		{
			description: "payload missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, payloadFlag)
			}),
			isValid: false,
		},
		{
			description: "payload is empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = ""
			}),
			isValid: false,
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
		expectedRequest loadbalancer.ApiUpdateLoadBalancerRequest
		isValid         bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
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
