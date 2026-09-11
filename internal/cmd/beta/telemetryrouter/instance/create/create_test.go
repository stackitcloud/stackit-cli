package create

import (
	"context"
	"testing"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	testRegion      = "eu01"
	testDisplayName = "my-telemetryrouter-instance"
	testDescription = "my instance description"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}
	testProjectId = uuid.NewString()
)

// Flags
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           testDisplayName,
		descriptionFlag:           testDescription,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

// Input Model
func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		DisplayName: new(testDisplayName),
		Description: new(testDescription),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureRequest(mods ...func(request *telemetryrouter.ApiCreateTelemetryRouterRequest)) telemetryrouter.ApiCreateTelemetryRouterRequest {
	request := testClient.DefaultAPI.CreateTelemetryRouter(testCtx, testProjectId, testRegion)
	request = request.CreateTelemetryRouterPayload(telemetryrouter.CreateTelemetryRouterPayload{
		DisplayName: testDisplayName,
		Description: new(testDescription),
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
		additionalFlagValues map[string][]string
		isValid              bool
		expectedModel        *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "optional flags omitted",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, descriptionFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
		},
		{
			description: "no values provided",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no filter-attribute at all still works, Filter stays nil",
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Filter = nil
			}),
		},
		{
			description: "single filter-attribute produces expected ConfigFilter",
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
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterAttributeFlag: {"key=http.method;level=bogus;matcher==;values=GET"},
			},
			isValid: false,
		},
		{
			description: "filter JSON produces expected ConfigFilter",
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
				// json.Unmarshal on the generated SDK types always yields a non-nil
				// (possibly empty) AdditionalProperties map, unlike BuildConfigFilter's
				// hand-constructed struct literals used by the other filter-attribute cases.
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
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterFlag: {`{"attributes": [`},
			},
			isValid: false,
		},
		{
			description: "filter and filter-attribute together is rejected",
			flagValues:  fixtureFlagValues(),
			additionalFlagValues: map[string][]string{
				filterFlag:          {`{"attributes": []}`},
				filterAttributeFlag: {testFilterAttributeRaw1},
			},
			isValid: false,
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
			description: "display name missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
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
		expectedRequest telemetryrouter.ApiCreateTelemetryRouterRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no optional values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
			expectedRequest: fixtureRequest().CreateTelemetryRouterPayload(telemetryrouter.CreateTelemetryRouterPayload{
				DisplayName: testDisplayName,
				Description: nil,
			}),
		},
		{
			description: "with filter attributes",
			model: fixtureInputModel(func(model *inputModel) {
				model.Filter = &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{testConfigFilterAttribute1, testConfigFilterAttribute2},
				}
			}),
			expectedRequest: fixtureRequest().CreateTelemetryRouterPayload(telemetryrouter.CreateTelemetryRouterPayload{
				DisplayName: testDisplayName,
				Description: new(testDescription),
				Filter: &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{testConfigFilterAttribute1, testConfigFilterAttribute2},
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(tt.expectedRequest, request,
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
			description: "default output",
			instance:    &telemetryrouter.TelemetryRouterResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			wantErr:     false,
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
