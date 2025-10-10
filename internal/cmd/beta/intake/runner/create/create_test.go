package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

// Define a unique key for the context to avoid collisions
type testCtxKey struct{}

const (
	testRegion             = "eu01"
	testDisplayName        = "testrunner"
	testMaxMessageSizeKiB  = int64(1024)
	testMaxMessagesPerHour = int64(10000)
	testDescription        = "This is a test runner"
	testLabelsString       = "env=test,team=dev"
)

var (
	// testCtx dummy context for testing purposes
	testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
	// testClient mock API client
	testClient    = &intake.APIClient{}
	testProjectId = uuid.NewString()

	testLabels = map[string]string{"env": "test", "team": "dev"}
)

// fixtureFlagValues generates a map of flag values for tests
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           testDisplayName,
		maxMessageSizeKiBFlag:     "1024",
		maxMessagesPerHourFlag:    "10000",
		descriptionFlag:           testDescription,
		labelFlag:                 testLabelsString,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

// fixtureInputModel generates an input model for tests
func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		DisplayName:        utils.Ptr(testDisplayName),
		MaxMessageSizeKiB:  utils.Ptr(testMaxMessageSizeKiB),
		MaxMessagesPerHour: utils.Ptr(testMaxMessagesPerHour),
		Description:        utils.Ptr(testDescription),
		Labels:             utils.Ptr(testLabels),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// fixtureCreatePayload generates a CreateIntakeRunnerPayload for tests
func fixtureCreatePayload(mods ...func(payload *intake.CreateIntakeRunnerPayload)) intake.CreateIntakeRunnerPayload {
	payload := intake.CreateIntakeRunnerPayload{
		DisplayName:        utils.Ptr(testDisplayName),
		MaxMessageSizeKiB:  utils.Ptr(testMaxMessageSizeKiB),
		MaxMessagesPerHour: utils.Ptr(testMaxMessagesPerHour),
		Description:        utils.Ptr(testDescription),
		Labels:             utils.Ptr(testLabels),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

// fixtureRequest generates an API request for tests
func fixtureRequest(mods ...func(request *intake.ApiCreateIntakeRunnerRequest)) intake.ApiCreateIntakeRunnerRequest {
	request := testClient.CreateIntakeRunner(testCtx, testProjectId, testRegion)
	request = request.CreateIntakeRunnerPayload(fixtureCreatePayload())
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
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "display name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "max message size missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, maxMessageSizeKiBFlag)
			}),
			isValid: false,
		},
		{
			description: "max messages per hour missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, maxMessagesPerHourFlag)
			}),
			isValid: false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				displayNameFlag:           testDisplayName,
				maxMessageSizeKiBFlag:     "1024",
				maxMessagesPerHourFlag:    "10000",
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.Labels = nil
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCreateCmd(&params.CmdParams{Printer: p})
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
		description     string
		model           *inputModel
		expectedRequest intake.ApiCreateIntakeRunnerRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no optionals",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.Labels = nil
			}),
			expectedRequest: fixtureRequest(func(request *intake.ApiCreateIntakeRunnerRequest) {
				*request = (*request).CreateIntakeRunnerPayload(fixtureCreatePayload(func(payload *intake.CreateIntakeRunnerPayload) {
					payload.Description = nil
					payload.Labels = nil
				}))
			}),
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

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		runnerId     string
		resp         *intake.IntakeRunnerResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "default output",
			args:    args{outputFormat: "default", projectLabel: "my-project", runnerId: "runner-id-123", resp: &intake.IntakeRunnerResponse{}},
			wantErr: false,
		},
		{
			name:    "json output",
			args:    args{outputFormat: print.JSONOutputFormat, resp: &intake.IntakeRunnerResponse{Id: utils.Ptr("runner-id-123")}},
			wantErr: false,
		},
		{
			name:    "yaml output",
			args:    args{outputFormat: print.YAMLOutputFormat, resp: &intake.IntakeRunnerResponse{Id: utils.Ptr("runner-id-123")}},
			wantErr: false,
		},
		{
			name:    "nil response",
			args:    args{outputFormat: print.JSONOutputFormat, resp: nil},
			wantErr: false,
		},
		{
			name:    "nil response - default output",
			args:    args{outputFormat: "default", resp: nil},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCreateCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
