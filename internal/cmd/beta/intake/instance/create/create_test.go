package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

// Define a unique key for the context to avoid collisions
type testCtxKey struct{}

const (
	testRegion = "eu01"

	testDisplayName            = "testintake"
	testDescription            = "This is a test intake"
	testLabelsString           = "env=test,team=dev"
	testCatalogURI             = "http://dremio.example.com"
	testCatalogWarehouse       = "my-warehouse"
	testCatalogNamespace       = "test-namespace"
	testCatalogTableName       = "test-table"
	testCatalogPartitioning    = "manual"
	testCatalogPartitionByFlag = "year,month"
	testCatalogAuthType        = "dremio"
	testDremioTokenEndpoint    = "https://auth.dremio.cloud/oauth/token" //nolint:gosec // false url
	testDremioToken            = "dremio-secret-token"
)

var (
	// testCtx dummy context for testing purposes
	testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
	// testClient mock API client
	testClient    = &intake.APIClient{}
	testProjectId = uuid.NewString()
	testRunnerId  = uuid.NewString()

	testLabels             = map[string]string{"env": "test", "team": "dev"}
	testCatalogPartitionBy = []string{"year", "month"}
)

// fixtureFlagValues generates a map of flag values for tests
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           testDisplayName,
		runnerIdFlag:              testRunnerId,
		descriptionFlag:           testDescription,
		labelsFlag:                testLabelsString,
		catalogURIFlag:            testCatalogURI,
		catalogWarehouseFlag:      testCatalogWarehouse,
		catalogNamespaceFlag:      testCatalogNamespace,
		catalogTableNameFlag:      testCatalogTableName,
		catalogPartitionByFlag:    testCatalogPartitionByFlag,
		catalogPartitioningFlag:   testCatalogPartitioning,
		catalogAuthTypeFlag:       testCatalogAuthType,
		dremioTokenEndpointFlag:   testDremioTokenEndpoint,
		dremioPatFlag:             testDremioToken,
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
		DisplayName:         utils.Ptr(testDisplayName),
		RunnerId:            utils.Ptr(testRunnerId),
		Description:         utils.Ptr(testDescription),
		Labels:              utils.Ptr(testLabels),
		CatalogURI:          utils.Ptr(testCatalogURI),
		CatalogWarehouse:    utils.Ptr(testCatalogWarehouse),
		CatalogNamespace:    utils.Ptr(testCatalogNamespace),
		CatalogTableName:    utils.Ptr(testCatalogTableName),
		CatalogPartitionBy:  utils.Ptr(testCatalogPartitionBy),
		CatalogPartitioning: utils.Ptr(testCatalogPartitioning),
		CatalogAuthType:     utils.Ptr(testCatalogAuthType),
		DremioTokenEndpoint: utils.Ptr(testDremioTokenEndpoint),
		DremioToken:         utils.Ptr(testDremioToken),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// fixtureCreatePayload generates a CreateIntakePayload for tests
func fixtureCreatePayload(mods ...func(payload *intake.CreateIntakePayload)) intake.CreateIntakePayload {
	authType := intake.CatalogAuthType(testCatalogAuthType)
	testPartitioningType := intake.PartitioningType(testCatalogPartitioning)
	payload := intake.CreateIntakePayload{
		DisplayName:    utils.Ptr(testDisplayName),
		IntakeRunnerId: utils.Ptr(testRunnerId),
		Description:    utils.Ptr(testDescription),
		Labels:         utils.Ptr(testLabels),
		Catalog: &intake.IntakeCatalog{
			Uri:          utils.Ptr(testCatalogURI),
			Warehouse:    utils.Ptr(testCatalogWarehouse),
			Namespace:    utils.Ptr(testCatalogNamespace),
			TableName:    utils.Ptr(testCatalogTableName),
			Partitioning: &testPartitioningType,
			PartitionBy:  utils.Ptr(testCatalogPartitionBy),
			Auth: &intake.CatalogAuth{
				Type: &authType,
				Dremio: &intake.DremioAuth{
					TokenEndpoint:       utils.Ptr(testDremioTokenEndpoint),
					PersonalAccessToken: utils.Ptr(testDremioToken),
				},
			},
		},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

// fixtureRequest generates an API request for tests
func fixtureRequest(mods ...func(request *intake.ApiCreateIntakeRequest)) intake.ApiCreateIntakeRequest {
	request := testClient.CreateIntake(testCtx, testProjectId, testRegion)
	request = request.CreateIntakePayload(fixtureCreatePayload())
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
			description: "runner-id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, runnerIdFlag)
			}),
			isValid: false,
		},
		{
			description: "catalog-uri missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, catalogURIFlag)
			}),
			isValid: false,
		},
		{
			description: "catalog-warehouse missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, catalogWarehouseFlag)
			}),
			isValid: false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				displayNameFlag:           testDisplayName,
				runnerIdFlag:              testRunnerId,
				catalogURIFlag:            testCatalogURI,
				catalogWarehouseFlag:      testCatalogWarehouse,
				catalogAuthTypeFlag:       testCatalogAuthType,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.Labels = nil
				model.CatalogNamespace = nil
				model.CatalogTableName = nil
				model.CatalogPartitioning = nil
				model.CatalogPartitionBy = nil
				model.DremioTokenEndpoint = nil
				model.DremioToken = nil
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, func(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
				return parseInput(p, cmd)
			}, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest intake.ApiCreateIntakeRequest
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
				model.CatalogNamespace = nil
				model.CatalogTableName = nil
				model.CatalogPartitioning = nil
				model.CatalogPartitionBy = nil
				model.CatalogAuthType = nil
				model.DremioTokenEndpoint = nil
				model.DremioToken = nil
			}),
			expectedRequest: fixtureRequest(func(request *intake.ApiCreateIntakeRequest) {
				*request = (*request).CreateIntakePayload(fixtureCreatePayload(func(payload *intake.CreateIntakePayload) {
					payload.Description = nil
					payload.Labels = nil
					payload.Catalog.Namespace = nil
					payload.Catalog.TableName = nil
					payload.Catalog.PartitionBy = nil
					payload.Catalog.Partitioning = nil
					payload.Catalog.Auth = nil
				}))
			}),
		},
		{
			description: "auth type none",
			model: fixtureInputModel(func(model *inputModel) {
				model.CatalogAuthType = utils.Ptr("none")
				model.DremioTokenEndpoint = nil
				model.DremioToken = nil
			}),
			expectedRequest: fixtureRequest(func(request *intake.ApiCreateIntakeRequest) {
				*request = (*request).CreateIntakePayload(fixtureCreatePayload(func(payload *intake.CreateIntakePayload) {
					authType := intake.CatalogAuthType("none")
					payload.Catalog.Auth.Type = &authType
					payload.Catalog.Auth.Dremio = nil
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
		model        *inputModel
		projectLabel string
		resp         *intake.IntakeResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default output",
			args: args{
				model:        fixtureInputModel(),
				projectLabel: "my-project",
				resp:         &intake.IntakeResponse{Id: utils.Ptr("intake-id-123")},
			},
			wantErr: false,
		},
		{
			name: "default output - async",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Async = true
				}),
				projectLabel: "my-project",
				resp:         &intake.IntakeResponse{Id: utils.Ptr("intake-id-123")},
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: &intake.IntakeResponse{Id: utils.Ptr("intake-id-123")},
			},
			wantErr: false,
		},
		{
			name: "nil response - default output",
			args: args{
				model: fixtureInputModel(),
				resp:  nil,
			},
			wantErr: false,
		},
		{
			name: "nil response - json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: nil,
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
