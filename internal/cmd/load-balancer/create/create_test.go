package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

const (
	testRegion = "eu02"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &loadbalancer.APIClient{}
	testProjectId = uuid.NewString()
	testRequestId = xRequestId
)

var testPayload = &loadbalancer.CreateLoadBalancerPayload{
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
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		payloadFlag: `{
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
  ]
}`,
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
		Payload: testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *loadbalancer.ApiCreateLoadBalancerRequest)) loadbalancer.ApiCreateLoadBalancerRequest {
	request := testClient.CreateLoadBalancer(testCtx, testProjectId, testRegion)
	request = request.CreateLoadBalancerPayload(*testPayload)
	request = request.XRequestID(testRequestId)
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
			description: "no flag values",
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
			description: "payload is missing",
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
		{
			description: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid:       false,
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

			err = cmd.ValidateFlagGroups()
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
			diff := cmp.Diff(*model, *tt.expectedModel,
				cmpopts.EquateComparable(testCtx),
			)
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
		expectedRequest loadbalancer.ApiCreateLoadBalancerRequest
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
