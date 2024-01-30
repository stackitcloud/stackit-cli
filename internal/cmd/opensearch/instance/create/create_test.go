package create

import (
	"context"
	"fmt"
	"testing"

	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/opensearch"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &opensearch.APIClient{}

type openSearchClientMocked struct {
	returnError       bool
	listOfferingsResp *opensearch.ListOfferingsResponse
}

func (c *openSearchClientMocked) CreateInstance(ctx context.Context, projectId string) opensearch.ApiCreateInstanceRequest {
	return testClient.CreateInstance(ctx, projectId)
}

func (c *openSearchClientMocked) ListOfferingsExecute(_ context.Context, _ string) (*opensearch.ListOfferingsResponse, error) {
	if c.returnError {
		return nil, fmt.Errorf("list flavors failed")
	}
	return c.listOfferingsResp, nil
}

var testProjectId = uuid.NewString()
var testPlanId = uuid.NewString()
var testMonitoringInstanceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:            testProjectId,
		instanceNameFlag:         "example-name",
		enableMonitoringFlag:     "true",
		graphiteFlag:             "example-graphite",
		metricsFrequencyFlag:     "100",
		metricsPrefixFlag:        "example-prefix",
		monitoringInstanceIdFlag: testMonitoringInstanceId,
		pluginFlag:               "example-plugin",
		sgwAclFlag:               "198.51.100.14/24",
		syslogFlag:               "example-syslog",
		planIdFlag:               testPlanId,
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
		},
		InstanceName:         utils.Ptr("example-name"),
		EnableMonitoring:     utils.Ptr(true),
		Graphite:             utils.Ptr("example-graphite"),
		MetricsFrequency:     utils.Ptr(int64(100)),
		MetricsPrefix:        utils.Ptr("example-prefix"),
		MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
		Plugin:               utils.Ptr([]string{"example-plugin"}),
		SgwAcl:               utils.Ptr([]string{"198.51.100.14/24"}),
		Syslog:               utils.Ptr([]string{"example-syslog"}),
		PlanId:               utils.Ptr(testPlanId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *opensearch.ApiCreateInstanceRequest)) opensearch.ApiCreateInstanceRequest {
	request := testClient.CreateInstance(testCtx, testProjectId)
	request = request.CreateInstancePayload(opensearch.CreateInstancePayload{
		InstanceName: utils.Ptr("example-name"),
		Parameters: &opensearch.InstanceParameters{
			EnableMonitoring:     utils.Ptr(true),
			Graphite:             utils.Ptr("example-graphite"),
			MetricsFrequency:     utils.Ptr(int64(100)),
			MetricsPrefix:        utils.Ptr("example-prefix"),
			MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
			Plugins:              utils.Ptr([]string{"example-plugin"}),
			SgwAcl:               utils.Ptr("198.51.100.14/24"),
			Syslog:               utils.Ptr([]string{"example-syslog"}),
		},
		PlanId: utils.Ptr(testPlanId),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
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
				projectIdFlag:    testProjectId,
				instanceNameFlag: "example-name",
				planIdFlag:       testPlanId,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				InstanceName: utils.Ptr("example-name"),
				PlanId:       utils.Ptr(testPlanId),
			},
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				projectIdFlag:        testProjectId,
				planIdFlag:           testPlanId,
				instanceNameFlag:     "",
				enableMonitoringFlag: "false",
				graphiteFlag:         "",
				metricsFrequencyFlag: "0",
				metricsPrefixFlag:    "",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				PlanId:           utils.Ptr(testPlanId),
				InstanceName:     utils.Ptr(""),
				EnableMonitoring: utils.Ptr(false),
				Graphite:         utils.Ptr(""),
				MetricsFrequency: utils.Ptr(int64(0)),
				MetricsPrefix:    utils.Ptr(""),
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
				model.Plugin = utils.Ptr(
					append(*model.Plugin, "example-plugin-1", "example-plugin-2"),
				)
			}),
		},
		{
			description:  "repeated syslog flags",
			flagValues:   fixtureFlagValues(),
			syslogValues: []string{"example-syslog-1", "example-syslog-2"},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Syslog = utils.Ptr(
					append(*model.Syslog, "example-syslog-1", "example-syslog-2"),
				)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd()
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

			for _, value := range tt.sgwAclValues {
				err := cmd.Flags().Set(sgwAclFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", sgwAclFlag, value, err)
				}
			}

			for _, value := range tt.pluginValues {
				err := cmd.Flags().Set(pluginFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", pluginFlag, value, err)
				}
			}

			for _, value := range tt.syslogValues {
				err := cmd.Flags().Set(syslogFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", syslogFlag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(cmd, "opensearch", "create")
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
		description       string
		model             *inputModel
		expectedRequest   opensearch.ApiCreateInstanceRequest
		getOfferingsFails bool
		getOfferingsResp  *opensearch.ListOfferingsResponse
		isValid           bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
			getOfferingsResp: &opensearch.ListOfferingsResponse{
				Offerings: &[]opensearch.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]opensearch.Plan{
							{
								Name: utils.Ptr("example-plan-name"),
								Id:   utils.Ptr(testPlanId),
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
			getOfferingsResp: &opensearch.ListOfferingsResponse{
				Offerings: &[]opensearch.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]opensearch.Plan{
							{
								Name: utils.Ptr("example-plan-name"),
								Id:   utils.Ptr(testPlanId),
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
			getOfferingsResp: &opensearch.ListOfferingsResponse{
				Offerings: &[]opensearch.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]opensearch.Plan{
							{
								Name: utils.Ptr("other-plan-name"),
								Id:   utils.Ptr(testPlanId),
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
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				PlanId: utils.Ptr(testPlanId),
			},
			getOfferingsResp: &opensearch.ListOfferingsResponse{
				Offerings: &[]opensearch.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]opensearch.Plan{
							{
								Name: utils.Ptr("example-plan-name"),
								Id:   utils.Ptr(testPlanId),
							},
						},
					},
				},
			},
			expectedRequest: testClient.CreateInstance(testCtx, testProjectId).
				CreateInstancePayload(opensearch.CreateInstancePayload{PlanId: utils.Ptr(testPlanId), Parameters: &opensearch.InstanceParameters{}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &openSearchClientMocked{
				returnError:       tt.getOfferingsFails,
				listOfferingsResp: tt.getOfferingsResp,
			}
			request, err := buildRequest(testCtx, "opensearch", tt.model, client)
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
