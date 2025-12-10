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
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
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
	testRunnerId  = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{testIntakeId}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		runnerIdFlag:              testRunnerId,
		displayNameFlag:           "new-display-name",
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
		RunnerId:    utils.Ptr(testRunnerId),
		DisplayName: utils.Ptr("new-display-name"),
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
				globalflags.RegionFlag:    testRegion,
				runnerIdFlag:              testRunnerId,
			},
			isValid: false,
		},
		{
			description: "update all fields",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[descriptionFlag] = "new description"
				flagValues[labelsFlag] = "env=prod,team=sre"
				flagValues[catalogURIFlag] = "new-uri"
				flagValues[catalogWarehouseFlag] = "new-warehouse"
				flagValues[catalogNamespaceFlag] = "new-namespace"
				flagValues[catalogTableNameFlag] = "new-table"
				flagValues[catalogAuthTypeFlag] = "dremio"
				flagValues[dremioTokenEndpointFlag] = "new-endpoint"
				flagValues[dremioPatFlag] = "new-pat"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("new description")
				model.Labels = utils.Ptr(map[string]string{"env": "prod", "team": "sre"})
				model.CatalogURI = utils.Ptr("new-uri")
				model.CatalogWarehouse = utils.Ptr("new-warehouse")
				model.CatalogNamespace = utils.Ptr("new-namespace")
				model.CatalogTableName = utils.Ptr("new-table")
				model.CatalogAuthType = utils.Ptr("dremio")
				model.DremioTokenEndpoint = utils.Ptr("new-endpoint")
				model.DremioToken = utils.Ptr("new-pat")
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
			description: "runner-id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, runnerIdFlag)
			}),
			isValid: false,
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
		description string
		model       *inputModel
		expectedReq intake.ApiUpdateIntakeRequest
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			expectedReq: testClient.UpdateIntake(testCtx, testProjectId, testRegion, testIntakeId).
				UpdateIntakePayload(intake.UpdateIntakePayload{
					IntakeRunnerId: utils.Ptr(testRunnerId),
					DisplayName:    utils.Ptr("new-display-name"),
				}),
		},
		{
			description: "update description and catalog uri",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Description = utils.Ptr("new-desc")
				model.CatalogURI = utils.Ptr("new-uri")
			}),
			expectedReq: testClient.UpdateIntake(testCtx, testProjectId, testRegion, testIntakeId).
				UpdateIntakePayload(intake.UpdateIntakePayload{
					IntakeRunnerId: utils.Ptr(testRunnerId),
					Description:    utils.Ptr("new-desc"),
					Catalog: &intake.IntakeCatalogPatch{
						Uri: utils.Ptr("new-uri"),
					},
				}),
		},
		{
			description: "update all fields",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = utils.Ptr("another-name")
				model.Description = utils.Ptr("final-desc")
				model.Labels = utils.Ptr(map[string]string{"a": "b"})
				model.CatalogURI = utils.Ptr("final-uri")
				model.CatalogWarehouse = utils.Ptr("final-warehouse")
				model.CatalogNamespace = utils.Ptr("final-namespace")
				model.CatalogTableName = utils.Ptr("final-table")
				model.CatalogAuthType = utils.Ptr("dremio")
				model.DremioTokenEndpoint = utils.Ptr("final-endpoint")
				model.DremioToken = utils.Ptr("final-token")
			}),
			expectedReq: testClient.UpdateIntake(testCtx, testProjectId, testRegion, testIntakeId).
				UpdateIntakePayload(intake.UpdateIntakePayload{
					IntakeRunnerId: utils.Ptr(testRunnerId),
					DisplayName:    utils.Ptr("another-name"),
					Description:    utils.Ptr("final-desc"),
					Labels:         utils.Ptr(map[string]string{"a": "b"}),
					Catalog: &intake.IntakeCatalogPatch{
						Uri:       utils.Ptr("final-uri"),
						Warehouse: utils.Ptr("final-warehouse"),
						Namespace: utils.Ptr("final-namespace"),
						TableName: utils.Ptr("final-table"),
						Auth: &intake.CatalogAuthPatch{
							Type: utils.Ptr(intake.CatalogAuthType("dremio")),
							Dremio: &intake.DremioAuthPatch{
								TokenEndpoint:       utils.Ptr("final-endpoint"),
								PersonalAccessToken: utils.Ptr("final-token"),
							},
						},
					},
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
		outputFormat string
		projectLabel string
		intakeId     string
		resp         *intake.IntakeResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "default output",
			args:    args{outputFormat: "default", projectLabel: "my-project", intakeId: "intake-id-123", resp: &intake.IntakeResponse{}},
			wantErr: false,
		},
		{
			name:    "json output",
			args:    args{outputFormat: print.JSONOutputFormat, resp: &intake.IntakeResponse{Id: utils.Ptr("intake-id-123")}},
			wantErr: false,
		},
		{
			name:    "yaml output",
			args:    args{outputFormat: print.YAMLOutputFormat, resp: &intake.IntakeResponse{Id: utils.Ptr("runner-id-123")}},
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
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: tt.args.outputFormat}}, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
