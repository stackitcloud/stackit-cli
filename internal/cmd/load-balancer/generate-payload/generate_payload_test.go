package generatepayload

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &loadbalancer.APIClient{}
var testProjectId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:        testProjectId,
		loadBalancerNameFlag: "example-name",
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
		LoadBalancerName: utils.Ptr("example-name"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *loadbalancer.ApiGetLoadBalancerRequest)) loadbalancer.ApiGetLoadBalancerRequest {
	request := testClient.GetLoadBalancer(testCtx, testProjectId, "example-name")
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
			isValid:     true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
			},
		},
		{
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, loadBalancerNameFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.LoadBalancerName = nil
			}),
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
		expectedRequest loadbalancer.ApiGetLoadBalancerRequest
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

func TestModifyListeners(t *testing.T) {
	tests := []struct {
		description string
		response    *loadbalancer.LoadBalancer
		expected    *[]loadbalancer.Listener
	}{
		{
			description: "base",
			response: &loadbalancer.LoadBalancer{
				Listeners: &[]loadbalancer.Listener{
					{
						DisplayName: utils.Ptr(""),
						Port:        utils.Ptr(int64(0)),
						Protocol:    utils.Ptr(""),
						Name:        utils.Ptr(""),
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
			},
			expected: &[]loadbalancer.Listener{
				{
					DisplayName: utils.Ptr(""),
					Port:        utils.Ptr(int64(0)),
					Protocol:    utils.Ptr(""),
					Name:        nil,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := modifyListener(tt.response)

			diff := cmp.Diff(output, tt.expected)
			if diff != "" {
				t.Errorf("expected output to be %+v, got %+v", tt.expected, output)
			}
		})
	}
}
