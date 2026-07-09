package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	redis "github.com/stackitcloud/stackit-sdk-go/services/redis/v2api"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &redis.APIClient{DefaultAPI: &redis.DefaultAPIService{}}
var testRegion = "eu01"

type mockSettings struct {
	returnError       bool
	listOfferingsResp *redis.ListOfferingsResponse
}

func newAPIMock(s mockSettings) redis.DefaultAPI {
	return &redis.DefaultAPIServiceMock{
		ListOfferingsExecuteMock: utils.Ptr(func(_ redis.ApiListOfferingsRequest) (*redis.ListOfferingsResponse, error) {
			if s.returnError {
				return nil, fmt.Errorf("list flavors failed")
			}
			return s.listOfferingsResp, nil
		}),
	}
}

var (
	testProjectId            = uuid.NewString()
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
		projectIdFlag:            testProjectId,
		globalflags.RegionFlag:   testRegion,
		enableMonitoringFlag:     "true",
		graphiteFlag:             "example-graphite",
		metricsFrequencyFlag:     "100",
		metricsPrefixFlag:        "example-prefix",
		monitoringInstanceIdFlag: testMonitoringInstanceId,
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
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		InstanceId:           testInstanceId,
		EnableMonitoring:     utils.Ptr(true),
		Graphite:             utils.Ptr("example-graphite"),
		MetricsFrequency:     utils.Ptr(int32(100)),
		MetricsPrefix:        utils.Ptr("example-prefix"),
		MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
		SgwAcl:               utils.Ptr([]string{"198.51.100.14/24"}),
		Syslog:               []string{"example-syslog"},
		PlanId:               utils.Ptr(testPlanId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *redis.ApiPartialUpdateInstanceRequest)) redis.ApiPartialUpdateInstanceRequest {
	request := testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.PartialUpdateInstancePayload(redis.PartialUpdateInstancePayload{
		Parameters: &redis.InstanceParameters{
			EnableMonitoring:     utils.Ptr(true),
			Graphite:             utils.Ptr("example-graphite"),
			MetricsFrequency:     utils.Ptr(int32(100)),
			MetricsPrefix:        utils.Ptr("example-prefix"),
			MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
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
				projectIdFlag: testProjectId,
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
				projectIdFlag:        testProjectId,
				planIdFlag:           testPlanId,
				enableMonitoringFlag: "false",
				graphiteFlag:         "",
				metricsFrequencyFlag: "0",
				metricsPrefixFlag:    "",
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
			description:  "repeated syslog flags",
			argValues:    fixtureArgValues(),
			flagValues:   fixtureFlagValues(),
			syslogValues: []string{"example-syslog-1", "example-syslog-2"},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Syslog = append(model.Syslog, "example-syslog-1", "example-syslog-2")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			cmd := NewCmd(params.CmdParams)
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

			for _, value := range tt.syslogValues {
				err := cmd.Flags().Set(syslogFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", syslogFlag, value, err)
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

			model, err := parseInput(params.Printer, cmd, tt.argValues)
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
		expectedRequest   redis.ApiPartialUpdateInstanceRequest
		getOfferingsFails bool
		listOfferingsResp *redis.ListOfferingsResponse
		isValid           bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
			listOfferingsResp: &redis.ListOfferingsResponse{
				Offerings: []redis.Offering{
					{
						Version: "example-version",
						Plans: []redis.Plan{
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
			listOfferingsResp: &redis.ListOfferingsResponse{
				Offerings: []redis.Offering{
					{
						Version: "example-version",
						Plans: []redis.Plan{
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
			listOfferingsResp: &redis.ListOfferingsResponse{
				Offerings: []redis.Offering{
					{
						Version: "example-version",
						Plans: []redis.Plan{
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
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				InstanceId: testInstanceId,
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testRegion, testInstanceId).
				PartialUpdateInstancePayload(redis.PartialUpdateInstancePayload{Parameters: &redis.InstanceParameters{}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := mockSettings{
				returnError:       tt.getOfferingsFails,
				listOfferingsResp: tt.listOfferingsResp,
			}
			request, err := buildRequest(testCtx, tt.model, newAPIMock(client))
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, redis.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
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
