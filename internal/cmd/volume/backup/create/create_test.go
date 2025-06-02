package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &iaas.APIClient{}
	testProjectId = uuid.NewString()
	testSourceId  = uuid.NewString()
	testName      = "my-backup"
	testLabels    = map[string]string{"key1": "value1"}
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		sourceIdFlag:              testSourceId,
		sourceTypeFlag:            "volume",
		nameFlag:                  testName,
		labelsFlag:                "key1=value1",
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
		SourceID:   testSourceId,
		SourceType: "volume",
		Name:       &testName,
		Labels:     testLabels,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateBackupRequest)) iaas.ApiCreateBackupRequest {
	request := testClient.CreateBackup(testCtx, testProjectId)
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
			description: "no source id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, sourceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "no source type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, sourceTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid source type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sourceTypeFlag] = "invalid"
			}),
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
			description: "project id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "only required flags",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
				delete(flagValues, labelsFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = nil
				model.Labels = make(map[string]string)
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
		expectedRequest iaas.ApiCreateBackupRequest
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

func TestOutputResult(t *testing.T) {
	backupId := "test-backup-id"

	type args struct {
		outputFormat string
		async        bool
		sourceLabel  string
		projectLabel string
		backup       *iaas.Backup
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty backup",
			args:    args{},
			wantErr: true,
		},
		{
			name: "minimal backup",
			args: args{
				backup: &iaas.Backup{
					Id: &backupId,
				},
				sourceLabel:  "test-source",
				projectLabel: "test-project",
			},
			wantErr: false,
		},
		{
			name: "async mode",
			args: args{
				backup: &iaas.Backup{
					Id: &backupId,
				},
				sourceLabel:  "test-source",
				projectLabel: "test-project",
				async:        true,
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				backup: &iaas.Backup{
					Id: &backupId,
				},
				outputFormat: print.JSONOutputFormat,
			},
			wantErr: false,
		},
		{
			name: "yaml output",
			args: args{
				backup: &iaas.Backup{
					Id: &backupId,
				},
				outputFormat: print.YAMLOutputFormat,
			},
			wantErr: false,
		},
	}

	p := print.NewPrinter()
	cmd := NewCmd(&params.CmdParams{Printer: p})
	p.Cmd = cmd

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.async, tt.args.sourceLabel, tt.args.projectLabel, tt.args.backup); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
