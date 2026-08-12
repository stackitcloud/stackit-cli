package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	projectIdFlag = globalflags.ProjectIdFlag
	testRegion    = "eu01"
)

type testCtxKey struct{}

var (
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient       = &valkey.APIClient{DefaultAPI: &valkey.DefaultAPIService{}}
	testProjectId    = uuid.NewString()
	testInstanceId   = uuid.NewString()
	testPlanId       = uuid.NewString()
	testMonitoringId = uuid.NewString()
)

type mockSettings struct {
	returnError       bool
	listOfferingsResp *valkey.ListOfferingsResponse
}

func newAPIMock(s mockSettings) valkey.DefaultAPI {
	return &valkey.DefaultAPIServiceMock{
		ListOfferingsExecuteMock: new(func(_ valkey.ApiListOfferingsRequest) (*valkey.ListOfferingsResponse, error) {
			if s.returnError {
				return nil, fmt.Errorf("list offerings failed")
			}
			return s.listOfferingsResp, nil
		}),
	}
}

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
		monitoringInstanceIdFlag: testMonitoringId,
		sgwAclFlag:               "198.51.100.14/24",
		syslogFlag:               "example-syslog",
		planIdFlag:               testPlanId,
		minReplicasToWriteFlag:   "2",
		replBacklogSizeFlag:      "10mb",
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
		PlanId:               new(testPlanId),
		EnableMonitoring:     new(true),
		Graphite:             new("example-graphite"),
		MetricsFrequency:     new(int32(100)),
		MetricsPrefix:        new("example-prefix"),
		MonitoringInstanceId: new(testMonitoringId),
		SgwAcl:               new([]string{"198.51.100.14/24"}),
		Syslog:               []string{"example-syslog"},
		MinReplicasToWrite:   new(int32(2)),
		ReplBacklogSize:      new("10mb"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *valkey.ApiPartialUpdateInstanceRequest)) valkey.ApiPartialUpdateInstanceRequest {
	request := testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.PartialUpdateInstancePayload(valkey.PartialUpdateInstancePayload{
		PlanId: new(testPlanId),
		Parameters: &valkey.InstanceParameters{
			EnableMonitoring:     new(true),
			Graphite:             new("example-graphite"),
			MetricsFrequency:     new(int32(100)),
			MetricsPrefix:        new("example-prefix"),
			MonitoringInstanceId: new(testMonitoringId),
			SgwAcl:               new("198.51.100.14/24"),
			Syslog:               []string{"example-syslog"},
			MinReplicasToWrite:   new(int32(2)),
			ReplBacklogSize:      new("10mb"),
		},
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
			description: "with plan name and version",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
				flagValues[versionFlag] = "8"
				delete(flagValues, planIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
				model.PlanName = "plan-name"
				model.Version = "8"
			}),
		},
		{
			description: "no plan selection",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, planIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
			}),
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
			description: "invalid plan: id and name",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
			}),
			isValid: false,
		},
		{
			description: "invalid plan: id and version",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[versionFlag] = "8"
			}),
			isValid: false,
		},
		{
			description: "invalid plan: name only",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
				delete(flagValues, planIdFlag)
			}),
			isValid: false,
		},
		{
			description: "empty update",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag:          testProjectId,
				globalflags.RegionFlag: testRegion,
			},
			isValid: false,
		},
		{
			description:  "repeated acl flags",
			argValues:    fixtureArgValues(),
			flagValues:   fixtureFlagValues(),
			sgwAclValues: []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:      true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SgwAcl = new(append(*model.SgwAcl, "198.51.100.14/24", "198.51.100.14/32"))
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
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				sgwAclFlag: tt.sgwAclValues,
				syslogFlag: tt.syslogValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description       string
		model             *inputModel
		expectedRequest   valkey.ApiPartialUpdateInstanceRequest
		getOfferingsFails bool
		getOfferingsResp  *valkey.ListOfferingsResponse
		isValid           bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
			getOfferingsResp: &valkey.ListOfferingsResponse{
				Offerings: []valkey.Offering{
					{
						Version: "example-version",
						Plans: []valkey.Plan{
							{
								Name: "example-plan-name",
								Id:   testPlanId,
							},
						},
					},
				},
			},
			isValid: true,
		},
		{
			description: "use plan name and version",
			model: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
				model.PlanName = "example-plan-name"
				model.Version = "example-version"
			}),
			expectedRequest: fixtureRequest(),
			getOfferingsResp: &valkey.ListOfferingsResponse{
				Offerings: []valkey.Offering{
					{
						Version: "example-version",
						Plans: []valkey.Plan{
							{
								Name: "example-plan-name",
								Id:   testPlanId,
							},
						},
					},
				},
			},
			isValid: true,
		},
		{
			description: "no plan selection",
			model: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
			}),
			expectedRequest: fixtureRequest(func(request *valkey.ApiPartialUpdateInstanceRequest) {
				*request = request.PartialUpdateInstancePayload(valkey.PartialUpdateInstancePayload{
					PlanId: nil,
					Parameters: &valkey.InstanceParameters{
						EnableMonitoring:     new(true),
						Graphite:             new("example-graphite"),
						MetricsFrequency:     new(int32(100)),
						MetricsPrefix:        new("example-prefix"),
						MonitoringInstanceId: new(testMonitoringId),
						SgwAcl:               new("198.51.100.14/24"),
						Syslog:               []string{"example-syslog"},
						MinReplicasToWrite:   new(int32(2)),
						ReplBacklogSize:      new("10mb"),
					},
				})
			}),
			isValid: true,
		},
		{
			description:       "get offerings fails",
			model:             fixtureInputModel(),
			getOfferingsFails: true,
			isValid:           false,
		},
		{
			description: "plan name not found",
			model: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
				model.PlanName = "example-plan-name"
				model.Version = "example-version"
			}),
			getOfferingsResp: &valkey.ListOfferingsResponse{
				Offerings: []valkey.Offering{
					{
						Version: "example-version",
						Plans: []valkey.Plan{
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
			description: "acl is joined into single string",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
				},
				InstanceId: testInstanceId,
				SgwAcl:     new([]string{"10.0.0.0/8", "192.168.1.0/24"}),
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testRegion, testInstanceId).
				PartialUpdateInstancePayload(valkey.PartialUpdateInstancePayload{
					Parameters: &valkey.InstanceParameters{
						SgwAcl: new("10.0.0.0/8,192.168.1.0/24"),
					},
				}),
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			apiMock := newAPIMock(mockSettings{
				returnError:       tt.getOfferingsFails,
				listOfferingsResp: tt.getOfferingsResp,
			})
			request, err := buildRequest(testCtx, tt.model, apiMock)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, valkey.DefaultAPIService{}),
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
