package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
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

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:            testProjectId,
		globalflags.RegionFlag:   testRegion,
		instanceNameFlag:         "example-name",
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
		InstanceName:         "example-name",
		EnableMonitoring:     new(true),
		Graphite:             new("example-graphite"),
		MetricsFrequency:     new(int32(100)),
		MetricsPrefix:        new("example-prefix"),
		MonitoringInstanceId: new(testMonitoringId),
		SgwAcl:               new([]string{"198.51.100.14/24"}),
		Syslog:               []string{"example-syslog"},
		PlanId:               new(testPlanId),
		MinReplicasToWrite:   new(int32(2)),
		ReplBacklogSize:      new("10mb"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *valkey.ApiCreateInstanceRequest)) valkey.ApiCreateInstanceRequest {
	request := testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(valkey.CreateInstancePayload{
		InstanceName: "example-name",
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
				flagValues[versionFlag] = "7"
				delete(flagValues, planIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
				model.PlanName = "plan-name"
				model.Version = "7"
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
					Verbosity: globalflags.VerbosityDefault,
				},
				InstanceName: "example-name",
				PlanId:       new(testPlanId),
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
					Verbosity: globalflags.VerbosityDefault,
				},
				PlanId:           new(testPlanId),
				InstanceName:     "",
				EnableMonitoring: new(false),
				Graphite:         new(""),
				MetricsFrequency: new(int32(0)),
				MetricsPrefix:    new(""),
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
				flagValues[versionFlag] = "7"
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
				model.SgwAcl = new(
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
				model.SgwAcl = new(
					append(*model.SgwAcl, "198.51.100.14/24", "198.51.100.14/32"),
				)
			}),
		},
		{
			description:  "repeated syslog flags",
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
		expectedRequest   valkey.ApiCreateInstanceRequest
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
		},
		{
			description: "get offerings fails",
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
			description: "required fields only",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				PlanId: new(testPlanId),
			},
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
			expectedRequest: testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion).
				CreateInstancePayload(valkey.CreateInstancePayload{PlanId: testPlanId, Parameters: &valkey.InstanceParameters{}}),
		},
		{
			description: "acl is joined into single string",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
				},
				PlanId: new(testPlanId),
				SgwAcl: new([]string{"10.0.0.0/8", "192.168.1.0/24"}),
			},
			expectedRequest: testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion).
				CreateInstancePayload(valkey.CreateInstancePayload{
					PlanId: testPlanId,
					Parameters: &valkey.InstanceParameters{
						SgwAcl: new("10.0.0.0/8,192.168.1.0/24"),
					},
				}),
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := mockSettings{
				returnError:       tt.getOfferingsFails,
				listOfferingsResp: tt.getOfferingsResp,
			}
			request, err := buildRequest(testCtx, tt.model, newAPIMock(client))
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

func TestOutputResult(t *testing.T) {
	type args struct {
		model        *inputModel
		projectLabel string
		instanceId   string
		resp         *valkey.CreateInstanceResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				model:        &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
				projectLabel: "",
				instanceId:   testMonitoringId,
				resp:         &valkey.CreateInstanceResponse{},
			},
			wantErr: false,
		},
		{
			name: "nil response",
			args: args{
				model:        &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
				projectLabel: "",
				instanceId:   testMonitoringId,
			},
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.projectLabel, tt.args.instanceId, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("TestOutputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
