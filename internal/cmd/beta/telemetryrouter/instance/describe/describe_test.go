package describe

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

const testRegion = "eu01"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
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
		InstanceID: testInstanceId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiGetTelemetryRouterRequest)) telemetryrouter.ApiGetTelemetryRouterRequest {
	request := testClient.DefaultAPI.GetTelemetryRouter(testCtx, testProjectId, testRegion, testInstanceId)
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
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argValues:     []string{testInstanceId},
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
			description: "project id missing",
			argValues:   []string{testInstanceId},
			flagValues:  fixtureFlagValues(func(m map[string]string) { delete(m, globalflags.ProjectIdFlag) }),
			isValid:     false,
		},
		{
			description: "project id invalid 1",
			argValues:   []string{testInstanceId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   []string{testInstanceId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "invalid instance id",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest telemetryrouter.ApiGetTelemetryRouterRequest
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
				cmpopts.EquateComparable(testCtx, telemetryrouter.DefaultAPIService{}),
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
		instance     *telemetryrouter.TelemetryRouterResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "default output",
			args:    args{outputFormat: "default", instance: &telemetryrouter.TelemetryRouterResponse{}},
			wantErr: false,
		},
		{
			name: "default output with filter",
			args: args{outputFormat: "default", instance: &telemetryrouter.TelemetryRouterResponse{
				Filter: &telemetryrouter.ConfigFilter{
					Attributes: []telemetryrouter.ConfigFilterAttributes{
						{
							Key:     "http.method",
							Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
							Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
							Values:  []string{"GET", "HEAD"},
						},
					},
				},
			}},
			wantErr: false,
		},
		{
			name:    "json output",
			args:    args{outputFormat: print.JSONOutputFormat, instance: &telemetryrouter.TelemetryRouterResponse{}},
			wantErr: false,
		},
		{
			name:    "yaml output",
			args:    args{outputFormat: print.YAMLOutputFormat, instance: &telemetryrouter.TelemetryRouterResponse{}},
			wantErr: false,
		},
		{
			name:    "nil instance",
			args:    args{instance: nil},
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.instance); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatFilterAttributes(t *testing.T) {
	tests := []struct {
		description string
		filter      *telemetryrouter.ConfigFilter
		want        string
	}{
		{
			description: "nil filter",
			filter:      nil,
			want:        "-",
		},
		{
			description: "empty attributes",
			filter:      &telemetryrouter.ConfigFilter{Attributes: []telemetryrouter.ConfigFilterAttributes{}},
			want:        "-",
		},
		{
			description: "single attribute",
			filter: &telemetryrouter.ConfigFilter{
				Attributes: []telemetryrouter.ConfigFilterAttributes{
					{
						Key:     "http.method",
						Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
						Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
						Values:  []string{"GET", "HEAD"},
					},
				},
			},
			want: "key=http.method;level=resource;matcher==;values=GET,HEAD",
		},
		{
			description: "multiple attributes joined by newline",
			filter: &telemetryrouter.ConfigFilter{
				Attributes: []telemetryrouter.ConfigFilterAttributes{
					{
						Key:     "http.method",
						Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
						Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
						Values:  []string{"GET", "HEAD"},
					},
					{
						Key:     "http.status_code",
						Level:   telemetryrouter.CONFIGFILTERLEVEL_LOG_RECORD,
						Matcher: telemetryrouter.CONFIGFILTERMATCHER_NOT_EQUAL,
						Values:  []string{"500", "503"},
					},
				},
			},
			want: "key=http.method;level=resource;matcher==;values=GET,HEAD\n" +
				"key=http.status_code;level=logRecord;matcher=!=;values=500,503",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := formatFilterAttributes(tt.filter)
			if got != tt.want {
				t.Fatalf("formatFilterAttributes() = %q, want %q", got, tt.want)
			}
		})
	}
}
