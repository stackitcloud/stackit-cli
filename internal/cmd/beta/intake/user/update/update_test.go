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

const (
	testRegion = "eu01"
)

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &intake.APIClient{}
	testProjectId = uuid.NewString()
	testIntakeId  = uuid.NewString()
	testUserId    = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{testUserId}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		intakeIdFlag:              testIntakeId,
		displayNameFlag:           "new-user-name",
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
		IntakeId:    testIntakeId,
		UserId:      testUserId,
		DisplayName: utils.Ptr("new-user-name"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			description: "no optional flags provided",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				intakeIdFlag:              testIntakeId,
			},
			isValid: false,
		},
		{
			description: "update all fields",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[passwordFlag] = "new-secret"
				flagValues[descriptionFlag] = "new description"
				flagValues[typeFlag] = "dead-letter"
				flagValues[labelFlag] = "env=prod,team=sre"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Password = utils.Ptr("new-secret")
				model.Description = utils.Ptr("new description")
				model.Type = utils.Ptr("dead-letter")
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
		{
			description: "intake id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, intakeIdFlag)
			}),
			isValid: false,
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
			diff := cmp.Diff(tt.expectedModel, model)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expectedReq intake.ApiUpdateIntakeUserRequest
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			expectedReq: testClient.UpdateIntakeUser(testCtx, testProjectId, testRegion, testIntakeId, testUserId).
				UpdateIntakeUserPayload(intake.UpdateIntakeUserPayload{
					DisplayName: utils.Ptr("new-user-name"),
				}),
		},
		{
			description: "update description and labels",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Description = utils.Ptr("new-desc")
				model.Labels = utils.Ptr(map[string]string{"key": "value"})
			}),
			expectedReq: testClient.UpdateIntakeUser(testCtx, testProjectId, testRegion, testIntakeId, testUserId).
				UpdateIntakeUserPayload(intake.UpdateIntakeUserPayload{
					Description: utils.Ptr("new-desc"),
					Labels:      utils.Ptr(map[string]string{"key": "value"}),
				}),
		},
		{
			description: "update all fields",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = utils.Ptr("another-name")
				model.Password = utils.Ptr("new-secret")
				model.Description = utils.Ptr("final-desc")
				model.Type = utils.Ptr("dead-letter")
				model.Labels = utils.Ptr(map[string]string{"a": "b"})
			}),
			expectedReq: testClient.UpdateIntakeUser(testCtx, testProjectId, testRegion, testIntakeId, testUserId).
				UpdateIntakeUserPayload(intake.UpdateIntakeUserPayload{
					DisplayName: utils.Ptr("another-name"),
					Password:    utils.Ptr("new-secret"),
					Description: utils.Ptr("final-desc"),
					Type:        (*intake.UserType)(utils.Ptr("dead-letter")),
					Labels:      utils.Ptr(map[string]string{"a": "b"}),
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(tt.expectedReq, request,
				cmp.AllowUnexported(request),
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
		model        *inputModel
		outputFormat string
		intakeId     string
		resp         *intake.IntakeUserResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "default output",
			args:    args{outputFormat: "default", intakeId: "intake-id-123", resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")}, model: fixtureInputModel()},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{outputFormat: print.JSONOutputFormat, resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")}, model: fixtureInputModel(func(model *inputModel) {
				model.OutputFormat = print.JSONOutputFormat
			})},
			wantErr: false,
		},
		{
			name: "yaml output",
			args: args{outputFormat: print.YAMLOutputFormat, resp: &intake.IntakeUserResponse{Id: utils.Ptr("user-id-123")}, model: fixtureInputModel(func(model *inputModel) {
				model.OutputFormat = print.YAMLOutputFormat
			})},
			wantErr: false,
		},
		{
			name:    "nil response - default output",
			args:    args{outputFormat: "default", intakeId: "intake-id-123", resp: nil, model: fixtureInputModel()},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p}) // p.Cmd is needed for the printer to have context.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.intakeId, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
