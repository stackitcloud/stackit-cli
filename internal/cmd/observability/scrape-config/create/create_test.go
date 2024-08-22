package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &observability.APIClient{}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

var testPayload = &observability.CreateScrapeConfigPayload{
	BasicAuth: &observability.CreateScrapeConfigPayloadBasicAuth{
		Username: utils.Ptr("username"),
		Password: utils.Ptr("password"),
	},
	BearerToken:     utils.Ptr("bearerToken"),
	HonorLabels:     utils.Ptr(true),
	HonorTimeStamps: utils.Ptr(true),
	MetricsPath:     utils.Ptr("/metrics"),
	JobName:         utils.Ptr("default-name"),
	MetricsRelabelConfigs: &[]observability.CreateScrapeConfigPayloadMetricsRelabelConfigsInner{
		{
			Action:       utils.Ptr("replace"),
			Modulus:      utils.Ptr(1.0),
			Regex:        utils.Ptr("regex"),
			Replacement:  utils.Ptr("replacement"),
			Separator:    utils.Ptr("separator"),
			SourceLabels: &[]string{"sourceLabel"},
			TargetLabel:  utils.Ptr("targetLabel"),
		},
	},
	Params: &map[string]interface{}{
		"key":  []interface{}{string("value1"), string("value2")},
		"key2": []interface{}{},
	},
	SampleLimit:    utils.Ptr(1.0),
	Scheme:         utils.Ptr("scheme"),
	ScrapeInterval: utils.Ptr("interval"),
	ScrapeTimeout:  utils.Ptr("timeout"),
	StaticConfigs: &[]observability.CreateScrapeConfigPayloadStaticConfigsInner{
		{
			Labels: &map[string]interface{}{
				"label":  "value",
				"label2": "value2",
			},
			Targets: &[]string{"target"},
		},
	},
	TlsConfig: &observability.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig{
		InsecureSkipVerify: utils.Ptr(true),
	},
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:  testProjectId,
		instanceIdFlag: testInstanceId,
		payloadFlag: `{
			"jobName": "default-name",
			"basicAuth": {
				"username": "username",
				"password": "password"
			},
			"bearerToken": "bearerToken",
			"honorLabels": true,
			"honorTimeStamps": true,
			"metricsPath": "/metrics",
			"metricsRelabelConfigs": [
				{
					"action": "replace",
					"modulus": 1.0,
					"regex": "regex",
					"replacement": "replacement",
					"separator": "separator",
					"sourceLabels": ["sourceLabel"],
					"targetLabel": "targetLabel"
				}
			],
			"params": {
				"key": ["value1", "value2"],
				"key2": []
			},
			"sampleLimit": 1.0,
			"scheme": "scheme",
			"scrapeInterval": "interval",
			"scrapeTimeout": "timeout",
			"staticConfigs": [
				{
					"labels": {
						"label": "value",
						"label2": "value2"
					},
					"targets": ["target"]
				}
			],
			"tlsConfig": {
				"insecureSkipVerify": true
			}	
		  }`,
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
		InstanceId: testInstanceId,
		Payload:    testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *observability.ApiCreateScrapeConfigRequest)) observability.ApiCreateScrapeConfigRequest {
	request := testClient.CreateScrapeConfig(testCtx, testInstanceId, testProjectId)
	request = request.CreateScrapeConfigPayload(*testPayload)
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
			description: "no flag values",
			flagValues:  map[string]string{},
			isValid:     false,
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
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "default config",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, payloadFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Payload = nil
			}),
		},
		{
			description: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid:       false,
			expectedModel: fixtureInputModel(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd(nil)
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

			model, err := parseInput(nil, cmd)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(*model, *tt.expectedModel,
				cmpopts.EquateComparable(testCtx),
			)
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
		expectedRequest observability.ApiCreateScrapeConfigRequest
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
