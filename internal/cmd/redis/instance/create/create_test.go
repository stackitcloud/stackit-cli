package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/redis"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &redis.APIClient{}

type redisClientMocked struct {
	returnError       bool
	listOfferingsResp *redis.ListOfferingsResponse
}

func (c *redisClientMocked) CreateInstance(ctx context.Context, projectId string) redis.ApiCreateInstanceRequest {
	return testClient.CreateInstance(ctx, projectId)
}

func (c *redisClientMocked) ListOfferingsExecute(_ context.Context, _ string) (*redis.ListOfferingsResponse, error) {
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
		},
		InstanceName:         utils.Ptr("example-name"),
		EnableMonitoring:     utils.Ptr(true),
		Graphite:             utils.Ptr("example-graphite"),
		MetricsFrequency:     utils.Ptr(int64(100)),
		MetricsPrefix:        utils.Ptr("example-prefix"),
		MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
		SgwAcl:               utils.Ptr([]string{"198.51.100.14/24"}),
		Syslog:               utils.Ptr([]string{"example-syslog"}),
		PlanId:               utils.Ptr(testPlanId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *redis.ApiCreateInstanceRequest)) redis.ApiCreateInstanceRequest {
	request := testClient.CreateInstance(testCtx, testProjectId)
	request = request.CreateInstancePayload(redis.CreateInstancePayload{
		InstanceName: utils.Ptr("example-name"),
		Parameters: &redis.InstanceParameters{
			EnableMonitoring:     utils.Ptr(true),
			Graphite:             utils.Ptr("example-graphite"),
			MetricsFrequency:     utils.Ptr(int64(100)),
			MetricsPrefix:        utils.Ptr("example-prefix"),
			MonitoringInstanceId: utils.Ptr(testMonitoringInstanceId),
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
					Verbosity: globalflags.VerbosityDefault,
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
					Verbosity: globalflags.VerbosityDefault,
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

			err = cmd.ValidateRequiredFlags()
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
		description       string
		model             *inputModel
		expectedRequest   redis.ApiCreateInstanceRequest
		getOfferingsFails bool
		getOfferingsResp  *redis.ListOfferingsResponse
		isValid           bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
			getOfferingsResp: &redis.ListOfferingsResponse{
				Offerings: &[]redis.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]redis.Plan{
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
			getOfferingsResp: &redis.ListOfferingsResponse{
				Offerings: &[]redis.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]redis.Plan{
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
			getOfferingsResp: &redis.ListOfferingsResponse{
				Offerings: &[]redis.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]redis.Plan{
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
					Verbosity: globalflags.VerbosityDefault,
				},
				PlanId: utils.Ptr(testPlanId),
			},
			getOfferingsResp: &redis.ListOfferingsResponse{
				Offerings: &[]redis.Offering{
					{
						Version: utils.Ptr("example-version"),
						Plans: &[]redis.Plan{
							{
								Name: utils.Ptr("example-plan-name"),
								Id:   utils.Ptr(testPlanId),
							},
						},
					},
				},
			},
			expectedRequest: testClient.CreateInstance(testCtx, testProjectId).
				CreateInstancePayload(redis.CreateInstancePayload{PlanId: utils.Ptr(testPlanId), Parameters: &redis.InstanceParameters{}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &redisClientMocked{
				returnError:       tt.getOfferingsFails,
				listOfferingsResp: tt.getOfferingsResp,
			}
			request, err := buildRequest(testCtx, tt.model, client)
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

func Test_outputResult(t *testing.T) {
	type args struct {
		model        *inputModel
		projectLabel string
		instanceId   string
		resp         *redis.CreateInstanceResponse
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
				instanceId:   testMonitoringInstanceId,
				resp:         &redis.CreateInstanceResponse{},
			},
			wantErr: false,
		},
		{
			name: "nil response",
			args: args{
				model:        &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
				projectLabel: "",
				instanceId:   testMonitoringInstanceId,
			},
			wantErr: true,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.projectLabel, tt.args.instanceId, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
