package update

import (
	"context"
	"testing"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

type testCtxKey struct{}

const (
	testRegion = "eu01"
)

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
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
		displayNameFlag:           "name",
		descriptionFlag:           "Example",
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
		InstanceID:  testInstanceId,
		DisplayName: new("name"),
		Description: new("Example"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiUpdateTelemetryRouterRequest)) telemetryrouter.ApiUpdateTelemetryRouterRequest {
	request := testClient.DefaultAPI.UpdateTelemetryRouter(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.UpdateTelemetryRouterPayload(telemetryrouter.UpdateTelemetryRouterPayload{
		DisplayName: new("name"),
		Description: new("Example"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

const (
	testFilterAttributeRaw1 = "key=http.method;level=resource;matcher==;values=GET,HEAD"
	testFilterAttributeRaw2 = "key=http.status;level=logRecord;matcher=!=;values=500"
)

var (
	testConfigFilterAttribute1 = telemetryrouter.ConfigFilterAttributes{
		Key:     "http.method",
		Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
		Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
		Values:  []string{"GET", "HEAD"},
	}
	testConfigFilterAttribute2 = telemetryrouter.ConfigFilterAttributes{
		Key:     "http.status",
		Level:   telemetryrouter.CONFIGFILTERLEVEL_LOG_RECORD,
		Matcher: telemetryrouter.CONFIGFILTERMATCHER_NOT_EQUAL,
		Values:  []string{"500"},
	}
)

// addAdditionalProperties returns a copy of attr with a non-nil AdditionalProperties map, as
// produced by json.Unmarshal on the generated SDK type (used for --filter JSON test fixtures).
func addAdditionalProperties(attr *telemetryrouter.ConfigFilterAttributes) telemetryrouter.ConfigFilterAttributes {
	withProperties := *attr
	withProperties.AdditionalProperties = map[string]interface{}{}
	return withProperties
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description          string
		argValues            []string
		flagValues           map[string]string
		primaryFlagValues    []string
		additionalFlagValues map[string][]string
		isValid              bool
		expectedModel        *inputModel
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
				globalflags.RegionFlag:    testRegion,
			},
			isValid: false,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				InstanceID: testInstanceId,
			},
		},
		{
			description: "update all fields",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				displayNameFlag:           "display-name",
				descriptionFlag:           "description",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				InstanceID:  testInstanceId,
				DisplayName: new("display-name"),
				Description: new("description"),
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
			description: "single filter-attribute produces expected ConfigFilter",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterAttributeFlag: {testFilterAttributeRaw1},
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Filter = &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{testConfigFilterAttribute1},
				}
			}),
		},
		{
			description: "multiple filter-attribute flags produce multiple attributes in order",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterAttributeFlag: {testFilterAttributeRaw1, testFilterAttributeRaw2},
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Filter = &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{testConfigFilterAttribute1, testConfigFilterAttribute2},
				}
			}),
		},
		{
			description: "invalid filter-attribute value produces a parseInput error",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterAttributeFlag: {"key=http.method;level=bogus;matcher==;values=GET"},
			},
			isValid: false,
		},
		{
			description: "filter JSON produces expected ConfigFilter",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterFlag: {
					`{"attributes": [` +
						`{"key": "http.method", "level": "resource", "matcher": "=", "values": ["GET", "HEAD"]},` +
						`{"key": "http.status", "level": "logRecord", "matcher": "!=", "values": ["500"]}` +
						`]}`,
				},
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Filter = &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{
						addAdditionalProperties(&testConfigFilterAttribute1),
						addAdditionalProperties(&testConfigFilterAttribute2),
					},
					AdditionalProperties: map[string]interface{}{},
				}
			}),
		},
		{
			description: "invalid filter JSON produces a parseInput error",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterFlag: {`{"attributes": [`},
			},
			isValid: false,
		},
		{
			description: "filter and filter-attribute together is rejected",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterFlag:          {`{"attributes": []}`},
				filterAttributeFlag: {testFilterAttributeRaw1},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.additionalFlagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest telemetryrouter.ApiUpdateTelemetryRouterRequest
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
				InstanceID: testInstanceId,
			},
			expectedRequest: testClient.DefaultAPI.UpdateTelemetryRouter(testCtx, testProjectId, testRegion, testInstanceId).
				UpdateTelemetryRouterPayload(telemetryrouter.UpdateTelemetryRouterPayload{}),
		},
		{
			description: "with filter attributes",
			model: fixtureInputModel(func(model *inputModel) {
				model.Filter = &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{testConfigFilterAttribute1, testConfigFilterAttribute2},
				}
			}),
			expectedRequest: fixtureRequest().UpdateTelemetryRouterPayload(telemetryrouter.UpdateTelemetryRouterPayload{
				DisplayName: new("name"),
				Description: new("Example"),
				Filter: &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{testConfigFilterAttribute1, testConfigFilterAttribute2},
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, telemetryrouter.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		instance    *telemetryrouter.TelemetryRouterResponse
		wantErr     bool
	}{
		{
			description: "nil response",
			instance:    nil,
			wantErr:     true,
		},
		{
			description: "default output",
			instance:    &telemetryrouter.TelemetryRouterResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			wantErr:     false,
		},
		{
			description: "model is nil",
			instance:    &telemetryrouter.TelemetryRouterResponse{},
			model:       nil,
			wantErr:     true,
		},
		{
			description: "global flag nil",
			instance:    &telemetryrouter.TelemetryRouterResponse{},
			model:       &inputModel{GlobalFlagModel: nil},
			wantErr:     true,
		},
		{
			description: "json output",
			instance:    &telemetryrouter.TelemetryRouterResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
			wantErr:     false,
		},
		{
			description: "yaml output",
			instance:    &telemetryrouter.TelemetryRouterResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.YAMLOutputFormat}},
			wantErr:     false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(params.Printer, tt.model, "label", tt.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
