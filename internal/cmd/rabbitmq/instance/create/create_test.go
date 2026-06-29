package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	rabbitmq "github.com/stackitcloud/stackit-sdk-go/services/rabbitmq/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &rabbitmq.APIClient{DefaultAPI: &rabbitmq.DefaultAPIService{}}

type mockSettings struct {
	returnError       bool
	listOfferingsResp *rabbitmq.ListOfferingsResponse
}

func newAPIMock(settings mockSettings) rabbitmq.DefaultAPI {
	return rabbitmq.DefaultAPIServiceMock{
		ListOfferingsExecuteMock: utils.Ptr(func(_ rabbitmq.ApiListOfferingsRequest) (*rabbitmq.ListOfferingsResponse, error) {
			if settings.returnError {
				return nil, fmt.Errorf("list flavors failed")
			}
			return settings.listOfferingsResp, nil
		}),
	}
}

var testProjectId = uuid.NewString()
var testPlanId = uuid.NewString()
var testMonitoringInstanceId = uuid.NewString()
var testInstanceName = utils.Ptr("instance")

const testRegion = "eu01"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceNameFlag:          "example-name",
		enableMonitoringFlag:      "true",
		graphiteFlag:              "example-graphite",
		metricsFrequencyFlag:      "100",
		metricsPrefixFlag:         "example-prefix",
		monitoringInstanceIdFlag:  testMonitoringInstanceId,
		pluginFlag:                "example-plugin",
		sgwAclFlag:                "198.51.100.14/24",
		syslogFlag:                "example-syslog",
		planIdFlag:                testPlanId,
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
		InstanceName:         utils.Ptr("example-name"),
		EnableMonitoring:     utils.Ptr(true),
		Graphite:             utils.Ptr("example-graphite"),
		MetricsFrequency:     utils.Ptr(int32(100)),
		MetricsPrefix:        utils.Ptr("example-prefix"),
		MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
		Plugin:               []rabbitmq.InstanceParametersPluginsInner{"example-plugin"},
		SgwAcl:               utils.Ptr([]string{"198.51.100.14/24"}),
		Syslog:               []string{"example-syslog"},
		PlanId:               utils.Ptr(testPlanId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *rabbitmq.ApiCreateInstanceRequest)) rabbitmq.ApiCreateInstanceRequest {
	request := testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(rabbitmq.CreateInstancePayload{
		InstanceName: "example-name",
		Parameters: &rabbitmq.InstanceParameters{
			EnableMonitoring:     utils.Ptr(true),
			Graphite:             utils.Ptr("example-graphite"),
			MetricsFrequency:     utils.Ptr(int32(100)),
			MetricsPrefix:        utils.Ptr("example-prefix"),
			MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
			Plugins:              []rabbitmq.InstanceParametersPluginsInner{"example-plugin"},
			SgwAcl:               utils.Ptr("198.51.100.14/24"),
			Syslog:               []string{"example-syslog"},
		},
		PlanId: testPlanId,
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
		sgwAclValues  []string
		pluginValues  []string
		syslogValues  []string
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
			description: "with plan name and version",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
				flagValues[versionFlag] = "6"
				delete(flagValues, planIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
				model.PlanName = "plan-name"
				model.Version = "6"
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				instanceNameFlag:          "example-name",
				planIdFlag:                testPlanId,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				InstanceName: utils.Ptr("example-name"),
				PlanId:       utils.Ptr(testPlanId),
			},
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				planIdFlag:                testPlanId,
				instanceNameFlag:          "",
				enableMonitoringFlag:      "false",
				graphiteFlag:              "",
				metricsFrequencyFlag:      "0",
				metricsPrefixFlag:         "",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				PlanId:           utils.Ptr(testPlanId),
				InstanceName:     utils.Ptr(""),
				EnableMonitoring: utils.Ptr(false),
				Graphite:         utils.Ptr(""),
				MetricsFrequency: utils.Ptr(int32(0)),
				MetricsPrefix:    utils.Ptr(""),
			},
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
			description: "invalid with plan ID, plan name and version",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
				flagValues[versionFlag] = "6"
			}),
			isValid: false,
		},
		{
			description: "invalid with plan ID and plan name",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
			}),
			isValid: false,
		},
		{
			description: "invalid with plan name only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
				delete(flagValues, planIdFlag)
			}),
			isValid: false,
		},
		{
			description:  "repeated acl flags",
			flagValues:   fixtureFlagValues(),
			sgwAclValues: []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SgwAcl = utils.Ptr(
					append(*model.SgwAcl, "198.51.100.14/24", "198.51.100.14/32"),
				)
			}),
		},
		{
			description:  "repeated acl flag with list value",
			flagValues:   fixtureFlagValues(),
			sgwAclValues: []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SgwAcl = utils.Ptr(
					append(*model.SgwAcl, "198.51.100.14/24", "198.51.100.14/32"),
				)
			}),
		},
		{
			description:  "repeated plugin flags",
			flagValues:   fixtureFlagValues(),
			pluginValues: []string{"example-plugin-1", "example-plugin-2"},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Plugin =
					append(model.Plugin, "example-plugin-1", "example-plugin-2")
			}),
		},
		{
			description:  "repeated syslog flags",
			flagValues:   fixtureFlagValues(),
			syslogValues: []string{"example-syslog-1", "example-syslog-2"},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Syslog =
					append(model.Syslog, "example-syslog-1", "example-syslog-2")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				sgwAclFlag: tt.sgwAclValues,
				syslogFlag: tt.syslogValues,
				pluginFlag: tt.pluginValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description       string
		model             *inputModel
		expectedRequest   rabbitmq.ApiCreateInstanceRequest
		getOfferingsFails bool
		getOfferingsResp  *rabbitmq.ListOfferingsResponse
		isValid           bool
	}{
		{
			//TODO: nil pointer due to plugins in model
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
			getOfferingsResp: &rabbitmq.ListOfferingsResponse{
				Offerings: []rabbitmq.Offering{
					{
						Version: "example-version",
						Plans: []rabbitmq.Plan{
							{
								Name: "example-plan-name",
								Id:   testPlanId,
							},
						},
					},
				},
			},
		},
		{
			description: "use plan name and version",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
					model.Version = "example-version"
				},
			),
			expectedRequest: fixtureRequest(),
			getOfferingsResp: &rabbitmq.ListOfferingsResponse{
				Offerings: []rabbitmq.Offering{
					{
						Version: "example-version",
						Plans: []rabbitmq.Plan{
							{
								Name: "example-plan-name",
								Id:   testPlanId,
							},
						},
					},
				},
			},
		},
		{
			description: "get offering fails",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
					model.Version = "example-version"
				},
			),
			getOfferingsFails: true,
			isValid:           false,
		},
		{
			description: "plan name not found",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
					model.Version = "example-version"
				},
			),
			getOfferingsResp: &rabbitmq.ListOfferingsResponse{
				Offerings: []rabbitmq.Offering{
					{
						Version: "example-version",
						Plans: []rabbitmq.Plan{
							{
								Name: "other-plan-name",
								Id:   testPlanId,
							},
						},
					},
				},
			},
			isValid: false,
		},
		{
			description: "required fields only",
			model: &inputModel{
				InstanceName: testInstanceName,
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				PlanId: utils.Ptr(testPlanId),
			},
			getOfferingsResp: &rabbitmq.ListOfferingsResponse{
				Offerings: []rabbitmq.Offering{
					{
						Version: "example-version",
						Plans: []rabbitmq.Plan{
							{
								Name: "example-plan-name",
								Id:   testPlanId,
							},
						},
					},
				},
			},
			expectedRequest: testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion).
				CreateInstancePayload(rabbitmq.CreateInstancePayload{InstanceName: *testInstanceName, PlanId: testPlanId, Parameters: &rabbitmq.InstanceParameters{}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := mockSettings{
				returnError:       tt.getOfferingsFails,
				listOfferingsResp: tt.getOfferingsResp,
			}
			request, err := buildRequest(testCtx, tt.model, newAPIMock(settings))
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, rabbitmq.DefaultAPIService{}),
				cmp.FilterPath(func(p cmp.Path) bool {
					return p.String() == "ApiService"
				}, cmp.Ignore()),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		model        inputModel
		projectLabel string
		instanceId   string
		resp         *rabbitmq.CreateInstanceResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				model: inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				resp: &rabbitmq.CreateInstanceResponse{},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				model: inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
			},
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, &tt.args.model, tt.args.projectLabel, tt.args.instanceId, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
