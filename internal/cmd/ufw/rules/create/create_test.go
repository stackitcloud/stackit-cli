package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	ufw "github.com/stackitcloud/stackit-sdk-go/services/ufw/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &ufw.APIClient{DefaultAPI: &ufw.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testRegion   = "eu01"
	testProduct  = "redis"
	testType     = "ACL"
	testSourceIp = "1.1.1.1/32"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		productFlag:               testProduct,
		typeFlag:                  testType,
		sourceIpFlag:              testSourceIp,
		instanceIdFlag:            testInstanceId,
		directionFlag:             "ingress",
		descriptionFlag:           "example-description",
		etherTypeFlag:             "IPv4",
		portRangeFlag:             "80-443",
		protocolFlag:              "TCP",
		offsetFlag:                "10",
		securityGroupIdFlag:       "example-sec-group",
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
		Product:         new(testProduct),
		Type:            new(testType),
		SourceIp:        new(testSourceIp),
		InstanceId:      new(testInstanceId),
		Direction:       new("ingress"),
		Description:     new("example-description"),
		EtherType:       new("IPv4"),
		PortRange:       new("80-443"),
		Protocol:        new("TCP"),
		Offset:          new(int32(10)),
		SecurityGroupId: new("example-sec-group"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *ufw.ApiCreateRuleRequest)) ufw.ApiCreateRuleRequest {
	request := testClient.DefaultAPI.CreateRule(testCtx, testProjectId, testRegion)
	request = request.CreateRulePayload(ufw.CreateRulePayload{
		Product:         testProduct,
		Type:            testType,
		SourceIP:        testSourceIp,
		InstanceId:      testInstanceId,
		Direction:       new("ingress"),
		Description:     new("example-description"),
		EtherType:       new("IPv4"),
		PortRange:       new("80-443"),
		Protocol:        new("TCP"),
		Offset:          new(int32(10)),
		SecurityGroupId: new("example-sec-group"),
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
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				productFlag:               testProduct,
				typeFlag:                  testType,
				sourceIpFlag:              testSourceIp,
				instanceIdFlag:            testInstanceId,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				Product:    new(testProduct),
				Type:       new(testType),
				SourceIp:   new(testSourceIp),
				InstanceId: new(testInstanceId),
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
			description: "region missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.RegionFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest ufw.ApiCreateRuleRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "required fields only",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				Product:    new(testProduct),
				Type:       new(testType),
				SourceIp:   new(testSourceIp),
				InstanceId: new(testInstanceId),
			},
			expectedRequest: testClient.DefaultAPI.CreateRule(testCtx, testProjectId, testRegion).
				CreateRulePayload(ufw.CreateRulePayload{
					Product:    testProduct,
					Type:       testType,
					SourceIP:   testSourceIp,
					InstanceId: testInstanceId,
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, ufw.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		async        bool
		projectLabel string
		rule         *ufw.CreateRuleResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "set empty response",
			args: args{
				rule: &ufw.CreateRuleResponse{},
			},
			wantErr: false,
		},
	}

	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.async, tt.args.projectLabel, tt.args.rule); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
