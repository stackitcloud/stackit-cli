package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	opensearch "github.com/stackitcloud/stackit-sdk-go/services/opensearch/v2api"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &opensearch.APIClient{DefaultAPI: &opensearch.DefaultAPIService{}}

type mockSettings struct {
	getOfferingsFails bool
	listOfferingsResp *opensearch.ListOfferingsResponse
}

func newAPIClientMock(c mockSettings) opensearch.DefaultAPI {
	return opensearch.DefaultAPIServiceMock{
		ListOfferingsExecuteMock: utils.Ptr(func(_ opensearch.ApiListOfferingsRequest) (*opensearch.ListOfferingsResponse, error) {
			if c.getOfferingsFails {
				return nil, fmt.Errorf("list flavors failed")
			}
			return c.listOfferingsResp, nil
		}),
	}
}

var (
	testProjectId            = uuid.NewString()
	testRegion               = "eu01"
	testInstanceId           = uuid.NewString()
	testPlanId               = uuid.NewString()
	testMonitoringInstanceId = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testInstanceId,
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
		enableMonitoringFlag:      "true",
		graphiteFlag:              "example-graphite",
		metricsFrequencyFlag:      "100",
		metricsPrefixFlag:         "example-prefix",
		monitoringInstanceIdFlag:  testMonitoringInstanceId,
		flagPlugins.Name():        string(opensearch.INSTANCEPARAMETERSPLUGINSINNER_REPOSITORY_AZURE),
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
		InstanceId:           testInstanceId,
		EnableMonitoring:     utils.Ptr(true),
		Graphite:             utils.Ptr("example-graphite"),
		MetricsFrequency:     utils.Ptr(int32(100)),
		MetricsPrefix:        utils.Ptr("example-prefix"),
		MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
		Plugin:               []opensearch.InstanceParametersPluginsInner{opensearch.INSTANCEPARAMETERSPLUGINSINNER_REPOSITORY_AZURE},
		SgwAcl:               utils.Ptr([]string{"198.51.100.14/24"}),
		Syslog:               []string{"example-syslog"},
		PlanId:               utils.Ptr(testPlanId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *opensearch.ApiPartialUpdateInstanceRequest)) opensearch.ApiPartialUpdateInstanceRequest {
	request := testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.PartialUpdateInstancePayload(opensearch.PartialUpdateInstancePayload{
		Parameters: &opensearch.InstanceParameters{
			EnableMonitoring:     utils.Ptr(true),
			Graphite:             utils.Ptr("example-graphite"),
			MetricsFrequency:     utils.Ptr(int32(100)),
			MetricsPrefix:        utils.Ptr("example-prefix"),
			MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
			Plugins:              []opensearch.InstanceParametersPluginsInner{opensearch.INSTANCEPARAMETERSPLUGINSINNER_REPOSITORY_AZURE},
			SgwAcl:               utils.Ptr("198.51.100.14/24"),
			Syslog:               []string{"example-syslog"},
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
			description: "required flags only (no values to update)",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
			},
			isValid: false,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				InstanceId: testInstanceId,
			},
		},
		{
			description: "zero values",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				planIdFlag:                testPlanId,
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
				InstanceId:       testInstanceId,
				PlanId:           utils.Ptr(testPlanId),
				EnableMonitoring: utils.Ptr(false),
				Graphite:         utils.Ptr(""),
				MetricsFrequency: utils.Ptr(int32(0)),
				MetricsPrefix:    utils.Ptr(""),
			},
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
			description: "instance id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "instance id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no acl flag",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, sgwAclFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SgwAcl = nil
			}),
		},
		{
			description:  "repeated acl flags",
			argValues:    fixtureArgValues(),
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
			argValues:    fixtureArgValues(),
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
			argValues:    fixtureArgValues(),
			flagValues:   fixtureFlagValues(),
			pluginValues: []string{string(opensearch.INSTANCEPARAMETERSPLUGINSINNER_REPOSITORY_AZURE), string(opensearch.INSTANCEPARAMETERSPLUGINSINNER_ANALYSIS_PHONETIC)},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Plugin =
					append(model.Plugin, opensearch.INSTANCEPARAMETERSPLUGINSINNER_REPOSITORY_AZURE, opensearch.INSTANCEPARAMETERSPLUGINSINNER_ANALYSIS_PHONETIC)
			}),
		},
		{
			description:  "repeated syslog flags",
			argValues:    fixtureArgValues(),
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
				sgwAclFlag:         tt.sgwAclValues,
				flagPlugins.Name(): tt.pluginValues,
				syslogFlag:         tt.syslogValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description        string
		model              *inputModel
		expectedRequest    opensearch.ApiPartialUpdateInstanceRequest
		mockClientSettings mockSettings
		isValid            bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
			mockClientSettings: mockSettings{
				listOfferingsResp: &opensearch.ListOfferingsResponse{
					Offerings: []opensearch.Offering{
						{
							Version: "example-version",
							Plans: []opensearch.Plan{
								{
									Name: "example-plan-name",
									Id:   testPlanId,
								},
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
			mockClientSettings: mockSettings{
				listOfferingsResp: &opensearch.ListOfferingsResponse{
					Offerings: []opensearch.Offering{
						{
							Version: "example-version",
							Plans: []opensearch.Plan{
								{
									Name: "example-plan-name",
									Id:   testPlanId,
								},
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
			mockClientSettings: mockSettings{
				getOfferingsFails: true,
			},
			isValid: false,
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
			mockClientSettings: mockSettings{
				listOfferingsResp: &opensearch.ListOfferingsResponse{
					Offerings: []opensearch.Offering{
						{
							Version: "example-version",
							Plans: []opensearch.Plan{
								{
									Name: "other-plan-name",
									Id:   testPlanId,
								},
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
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				InstanceId: testInstanceId,
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testRegion, testInstanceId).
				PartialUpdateInstancePayload(opensearch.PartialUpdateInstancePayload{Parameters: &opensearch.InstanceParameters{}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, newAPIClientMock(tt.mockClientSettings))
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmpopts.IgnoreFields(tt.expectedRequest, "ApiService"),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
