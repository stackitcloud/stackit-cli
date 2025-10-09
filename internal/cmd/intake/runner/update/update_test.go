package update

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

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &intake.APIClient{}
	testProjectId = uuid.NewString()
	testRunnerId  = uuid.NewString()
	testRegion    = "eu01"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testRunnerId,
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
		displayNameFlag:           "new-runner-name",
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
		RunnerId:    testRunnerId,
		DisplayName: utils.Ptr("new-runner-name"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *intake.ApiUpdateIntakeRunnerRequest)) intake.ApiUpdateIntakeRunnerRequest {
	request := testClient.UpdateIntakeRunner(testCtx, testProjectId, testRegion, testRunnerId)
	payload := intake.UpdateIntakeRunnerPayload{
		DisplayName: utils.Ptr("new-runner-name"),
	}
	request = request.UpdateIntakeRunnerPayload(payload)
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no update flags provided",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
			},
			isValid: false,
		},
		{
			description: "update all fields",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[maxMessageSizeKiBFlag] = "2048"
				flagValues[maxMessagesPerHourFlag] = "10000"
				flagValues[descriptionFlag] = "new description"
				flagValues[labelFlag] = "env=prod,team=sre"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.MaxMessageSizeKiB = utils.Ptr(int64(2048))
				model.MaxMessagesPerHour = utils.Ptr(int64(10000))
				model.Description = utils.Ptr("new description")
				model.Labels = utils.Ptr(map[string]string{"env": "prod", "team": "sre"})
			}),
		},
		{
			description: "no args",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewUpdateCmd(&params.CmdParams{Printer: p})
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

			model, err := parseInput(p, cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
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
		expectedRequest intake.ApiUpdateIntakeRunnerRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "update description and labels",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Description = utils.Ptr("new-desc")
				model.Labels = utils.Ptr(map[string]string{"key": "value"})
			}),
			expectedRequest: fixtureRequest(func(request *intake.ApiUpdateIntakeRunnerRequest) {
				payload := intake.UpdateIntakeRunnerPayload{
					Description: utils.Ptr("new-desc"),
					Labels:      utils.Ptr(map[string]string{"key": "value"}),
				}
				*request = (*request).UpdateIntakeRunnerPayload(payload)
			}),
		},
		{
			description: "update all fields",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = utils.Ptr("another-name")
				model.MaxMessageSizeKiB = utils.Ptr(int64(4096))
				model.MaxMessagesPerHour = utils.Ptr(int64(20000))
				model.Description = utils.Ptr("final-desc")
				model.Labels = utils.Ptr(map[string]string{"a": "b"})
			}),
			expectedRequest: fixtureRequest(func(request *intake.ApiUpdateIntakeRunnerRequest) {
				payload := intake.UpdateIntakeRunnerPayload{
					DisplayName:        utils.Ptr("another-name"),
					MaxMessageSizeKiB:  utils.Ptr(int64(4096)),
					MaxMessagesPerHour: utils.Ptr(int64(20000)),
					Description:        utils.Ptr("final-desc"),
					Labels:             utils.Ptr(map[string]string{"a": "b"}),
				}
				*request = (*request).UpdateIntakeRunnerPayload(payload)
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
	p.Cmd = NewUpdateCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
