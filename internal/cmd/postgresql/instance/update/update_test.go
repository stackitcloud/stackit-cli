package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &postgresql.APIClient{}

type postgreSQLClientMocked struct {
	returnError      bool
	getOfferingsResp *postgresql.OfferingList
}

func (c *postgreSQLClientMocked) CreateInstance(_ context.Context, _ string) postgresql.ApiCreateInstanceRequest {
	return postgresql.ApiCreateInstanceRequest{}
}

func (c *postgreSQLClientMocked) UpdateInstance(_ context.Context, _, _ string) postgresql.ApiUpdateInstanceRequest {
	return postgresql.ApiUpdateInstanceRequest{}
}

func (c *postgreSQLClientMocked) GetOfferingsExecute(_ context.Context, _ string) (*postgresql.OfferingList, error) {
	if c.returnError {
		return nil, fmt.Errorf("get flavors failed")
	}
	return c.getOfferingsResp, nil
}

var (
	testProjectId                    = uuid.NewString()
	testInstanceId                   = uuid.NewString()
	currentPlanId                    = uuid.NewString()
	currentMonitoringInstanceIdValue = uuid.NewString()
	updatedPlanIdValue               = uuid.NewString()
	updatedMonitoringInstanceIdValue = uuid.NewString()
)

const (
	currentEnableMonitoringValue = false
	currentGraphiteValue         = "example-graphite"
	currentMetricsFrequencyValue = int64(100)
	currentMetricsPrefixValue    = "example-prefix"
	currentPluginValue           = "example-plugin"
	currentSgwAclValue           = "198.51.100.14/24"
	currentSyslogValue           = "example-syslog"
	updatedEnableMonitoringValue = true
	updatedGraphiteValue         = "example-graphite-updated"
	updatedMetricsFrequencyValue = int64(101)
	updatedMetricsPrefixValue    = "example-prefix-updated"
	updatedPluginValue           = "example-plugin-updated"
	updatedSgwAclValue           = "0.0.0.0/0"
	updatedSyslogValue           = "example-syslog-updated"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:            testProjectId,
		instanceIdFlag:           testInstanceId,
		enableMonitoringFlag:     "true",
		graphiteFlag:             updatedGraphiteValue,
		metricsFrequencyFlag:     "101",
		metricsPrefixFlag:        updatedMetricsPrefixValue,
		monitoringInstanceIdFlag: updatedMonitoringInstanceIdValue,
		pluginFlag:               updatedPluginValue,
		sgwAclFlag:               updatedSgwAclValue,
		syslogFlag:               updatedSyslogValue,
		planIdFlag:               updatedPlanIdValue,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagModel(mods ...func(model *flagModel)) *flagModel {
	model := &flagModel{
		ProjectId:            testProjectId,
		InstanceId:           testInstanceId,
		EnableMonitoring:     utils.Ptr(updatedEnableMonitoringValue),
		Graphite:             utils.Ptr(updatedGraphiteValue),
		MetricsFrequency:     utils.Ptr(updatedMetricsFrequencyValue),
		MetricsPrefix:        utils.Ptr(updatedMetricsPrefixValue),
		MonitoringInstanceId: utils.Ptr(updatedMonitoringInstanceIdValue),
		Plugin:               utils.Ptr([]string{updatedPluginValue}),
		SgwAcl:               utils.Ptr([]string{updatedSgwAclValue}),
		Syslog:               utils.Ptr([]string{updatedSyslogValue}),
		PlanId:               utils.Ptr(updatedPlanIdValue),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureInstance(mods ...func(instance *postgresql.Instance)) *postgresql.Instance {
	instance := &postgresql.Instance{
		InstanceId: utils.Ptr(testInstanceId),
		PlanId:     utils.Ptr(currentPlanId),
		Name:       utils.Ptr("example-name"),
		Parameters: utils.Ptr(map[string]interface{}{
			"enable_monitoring":      currentEnableMonitoringValue,
			"graphite":               currentGraphiteValue,
			"metrics_frequency":      currentMetricsFrequencyValue,
			"metrics_prefix":         currentMetricsPrefixValue,
			"monitoring_instance_id": currentMonitoringInstanceIdValue,
			"plugins":                []string{currentPluginValue},
			"sgw_acl":                currentSgwAclValue,
			"syslog":                 []string{currentSyslogValue},
		}),
	}
	for _, mod := range mods {
		mod(instance)
	}
	return instance
}

func fixtureRequest(mods ...func(request *postgresql.ApiUpdateInstanceRequest)) postgresql.ApiUpdateInstanceRequest {
	request := testClient.UpdateInstance(testCtx, testProjectId, testInstanceId)
	request = request.UpdateInstancePayload(postgresql.UpdateInstancePayload{
		Parameters: &postgresql.InstanceParameters{
			EnableMonitoring:     utils.Ptr(updatedEnableMonitoringValue),
			Graphite:             utils.Ptr(updatedGraphiteValue),
			MetricsFrequency:     utils.Ptr(updatedMetricsFrequencyValue),
			MetricsPrefix:        utils.Ptr(updatedMetricsPrefixValue),
			MonitoringInstanceId: utils.Ptr(updatedMonitoringInstanceIdValue),
			Plugins:              utils.Ptr([]string{updatedPluginValue}),
			SgwAcl:               utils.Ptr(updatedSgwAclValue),
			Syslog:               utils.Ptr([]string{updatedSyslogValue}),
		},
		PlanId: utils.Ptr(updatedPlanIdValue),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		sgwAclValues  []string
		pluginValues  []string
		syslogValues  []string
		isValid       bool
		expectedModel *flagModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureFlagModel(),
		},
		{
			description: "base with plan_name and version",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, planIdFlag)
					flagValues[planNameFlag] = "example-plan-name"
					flagValues[versionFlag] = "example-version"
				},
			),
			isValid: true,
			expectedModel: fixtureFlagModel(
				func(model *flagModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
					model.Version = "example-version"
				},
			),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				projectIdFlag:  testProjectId,
				instanceIdFlag: testInstanceId,
			},
			isValid: true,
			expectedModel: &flagModel{
				ProjectId:  testProjectId,
				InstanceId: testInstanceId,
			},
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				projectIdFlag:        testProjectId,
				planIdFlag:           currentPlanId,
				instanceIdFlag:       testInstanceId,
				enableMonitoringFlag: "false",
				graphiteFlag:         "",
				metricsFrequencyFlag: "0",
				metricsPrefixFlag:    "",
			},
			isValid: true,
			expectedModel: &flagModel{
				ProjectId:        testProjectId,
				InstanceId:       testInstanceId,
				PlanId:           utils.Ptr(currentPlanId),
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
			description:  "repeated acl flags",
			flagValues:   fixtureFlagValues(),
			sgwAclValues: []string{currentSgwAclValue, "198.51.100.14/32"},
			isValid:      true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.SgwAcl = utils.Ptr(
					append(*model.SgwAcl, currentSgwAclValue, "198.51.100.14/32"),
				)
			}),
		},
		{
			description:  "repeated acl flag with list value",
			flagValues:   fixtureFlagValues(),
			sgwAclValues: []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:      true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.SgwAcl = utils.Ptr(
					append(*model.SgwAcl, currentSgwAclValue, "198.51.100.14/32"),
				)
			}),
		},
		{
			description:  "repeated plugin flags",
			flagValues:   fixtureFlagValues(),
			pluginValues: []string{"example-plugin-1", "example-plugin-2"},
			isValid:      true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
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
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.Syslog = utils.Ptr(
					append(*model.Syslog, "example-syslog-1", "example-syslog-2"),
				)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Flag defined in root command
			err := testutils.ConfigureBindUUIDFlag(cmd, projectIdFlag, config.ProjectIdKey)
			if err != nil {
				t.Fatalf("configure global flag --%s: %v", projectIdFlag, err)
			}

			configureFlags(cmd)

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

			model, err := parseFlags(cmd)
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
		model             *flagModel
		instance          *postgresql.Instance
		expectedRequest   postgresql.ApiUpdateInstanceRequest
		getOfferingsFails bool
		getOfferingsResp  *postgresql.OfferingList
		isValid           bool
	}{
		{
			description:     "base",
			instance:        fixtureInstance(),
			model:           fixtureFlagModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "do not update any field",
			instance:    fixtureInstance(),
			model: fixtureFlagModel(
				func(model *flagModel) {
					model.PlanId = nil
					model.PlanName = ""
					model.Version = ""
					model.EnableMonitoring = nil
					model.Graphite = nil
					model.MetricsFrequency = nil
					model.MetricsPrefix = nil
					model.MonitoringInstanceId = nil
					model.Plugin = nil
					model.SgwAcl = nil
					model.Syslog = nil
				},
			),
			expectedRequest: testClient.UpdateInstance(testCtx, testProjectId, testInstanceId).
				UpdateInstancePayload(postgresql.UpdateInstancePayload{
					Parameters: &postgresql.InstanceParameters{
						EnableMonitoring:     utils.Ptr(currentEnableMonitoringValue),
						Graphite:             utils.Ptr(currentGraphiteValue),
						MetricsFrequency:     utils.Ptr(currentMetricsFrequencyValue),
						MetricsPrefix:        utils.Ptr(currentMetricsPrefixValue),
						MonitoringInstanceId: utils.Ptr(currentMonitoringInstanceIdValue),
						Plugins:              utils.Ptr([]string{currentPluginValue}),
						SgwAcl:               utils.Ptr(currentSgwAclValue),
						Syslog:               utils.Ptr([]string{currentSyslogValue}),
					},
					PlanId: utils.Ptr(currentPlanId),
				}),
		},
		{
			description: "use plan name and version",
			model: fixtureFlagModel(
				func(model *flagModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
					model.Version = "example-version"
				},
			),
			instance:        fixtureInstance(),
			expectedRequest: fixtureRequest(),
			getOfferingsResp: &postgresql.OfferingList{
				Offerings: &[]postgresql.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]postgresql.Plan{
							{
								Name: utils.Ptr("example-plan-name"),
								Id:   utils.Ptr(updatedPlanIdValue),
							},
						},
					},
				},
			},
		},
		{
			description: "get offering fails",
			model: fixtureFlagModel(
				func(model *flagModel) {
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
			model: fixtureFlagModel(
				func(model *flagModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
					model.Version = "example-version"
				},
			),
			getOfferingsResp: &postgresql.OfferingList{
				Offerings: &[]postgresql.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]postgresql.Plan{
							{
								Name: utils.Ptr("other-plan-name"),
								Id:   utils.Ptr(currentPlanId),
							},
						},
					},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &postgreSQLClientMocked{
				returnError:      tt.getOfferingsFails,
				getOfferingsResp: tt.getOfferingsResp,
			}
			request, err := buildRequest(testCtx, tt.instance, tt.model, client)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.IgnoreFields(postgresql.ApiUpdateInstanceRequest{}, "apiService", "ctx", "projectId", "instanceId"),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
