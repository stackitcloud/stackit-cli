package options

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testInstanceId = uuid.NewString()

type mockSettings struct {
	listFlavorsFails           bool
	listVersionsFails          bool
	listStoragesFails          bool
	listUserRolesFails         bool
	listDBCollationsFails      bool
	listDBCompatibilitiesFails bool

	listFlavorsCalled           bool
	listVersionsCalled          bool
	listStoragesCalled          bool
	listUserRolesCalled         bool
	listDBCollationsCalled      bool
	listDBCompatibilitiesCalled bool
}

func newAPIMock(s *mockSettings) sqlserverflex.DefaultAPI {
	return &sqlserverflex.DefaultAPIServiceMock{
		ListFlavorsExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListFlavorsRequest) (*sqlserverflex.ListFlavorsResponse, error) {
			s.listFlavorsCalled = true
			if s.listFlavorsFails {
				return nil, fmt.Errorf("list flavors failed")
			}
			return utils.Ptr(sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{},
			}), nil
		}),
		ListVersionsExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListVersionsRequest) (*sqlserverflex.ListVersionsResponse, error) {
			s.listVersionsCalled = true
			if s.listVersionsFails {
				return nil, fmt.Errorf("list versions failed")
			}
			return utils.Ptr(sqlserverflex.ListVersionsResponse{
				Versions: []sqlserverflex.Version{},
			}), nil
		}),
		ListStoragesExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListStoragesRequest) (*sqlserverflex.ListStoragesResponse, error) {
			s.listStoragesCalled = true
			if s.listStoragesFails {
				return nil, fmt.Errorf("list storages failed")
			}
			return utils.Ptr(sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(10),
					Max: int32(100),
				},
			}), nil
		}),
		ListRolesExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListRolesRequest) (*sqlserverflex.ListRolesResponse, error) {
			s.listUserRolesCalled = true
			if s.listUserRolesFails {
				return nil, fmt.Errorf("list roles failed")
			}
			return utils.Ptr(sqlserverflex.ListRolesResponse{
				Roles: []string{},
			}), nil
		}),
		ListCollationsExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListCollationsRequest) (*sqlserverflex.ListCollationsResponse, error) {
			s.listDBCollationsCalled = true
			if s.listDBCollationsFails {
				return nil, fmt.Errorf("list collations failed")
			}
			return utils.Ptr(sqlserverflex.ListCollationsResponse{
				Collations: []sqlserverflex.DatabaseGetcollation{},
			}), nil
		}),
		ListCompatibilitiesExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListCompatibilitiesRequest) (*sqlserverflex.ListCompatibilityResponse, error) {
			s.listDBCompatibilitiesCalled = true
			if s.listDBCompatibilitiesFails {
				return nil, fmt.Errorf("list compatibilities failed")
			}
			return utils.Ptr(sqlserverflex.ListCompatibilityResponse{
				Compatibilities: []sqlserverflex.DatabaseGetcompatibility{},
			}), nil
		}),
	}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		flavorsFlag:           "true",
		versionsFlag:          "true",
		storagesFlag:          "true",
		userRolesFlag:         "true",
		dbCollationsFlag:      "true",
		dbCompatibilitiesFlag: "true",
		flavorIdFlag:          "2.4",
		instanceIdFlag:        testInstanceId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModelAllFalse(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel:   &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		Flavors:           false,
		Versions:          false,
		Storages:          false,
		UserRoles:         false,
		DBCollations:      false,
		DBCompatibilities: false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureInputModelAllTrue(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel:   &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		Flavors:           true,
		Versions:          true,
		Storages:          true,
		UserRoles:         true,
		DBCollations:      true,
		DBCompatibilities: true,
		FlavorId:          utils.Ptr("2.4"),
		InstanceId:        utils.Ptr(testInstanceId),
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
			description:   "all values",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModelAllTrue(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "some values 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[storagesFlag] = "false"
				delete(flagValues, flavorIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModelAllTrue(func(model *inputModel) {
				model.Storages = false
				model.FlavorId = nil
			}),
		},
		{
			description: "some values 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, flavorsFlag)
				delete(flagValues, versionsFlag)
				delete(flagValues, userRolesFlag)
				flagValues[storagesFlag] = "true"
				flagValues[flavorIdFlag] = "2.4"
			}),
			isValid: true,
			expectedModel: fixtureInputModelAllTrue(func(model *inputModel) {
				model.Flavors = false
				model.Versions = false
				model.UserRoles = false
			}),
		},
		{
			description: "storages without flavor-id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, flavorIdFlag)
			}),
			isValid: false,
		},
		{
			description: "flavor-id without storage",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, storagesFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModelAllTrue(func(model *inputModel) {
				model.Storages = false
			}),
		},
		{
			description: "user roles without instance id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
				delete(flagValues, dbCollationsFlag)
				delete(flagValues, dbCompatibilitiesFlag)
			}),
			isValid: false,
		},
		{
			description: "db collations without instance id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
				delete(flagValues, userRolesFlag)
				delete(flagValues, dbCompatibilitiesFlag)
			}),
			isValid: false,
		},
		{
			description: "db compatibilities without instance id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
				delete(flagValues, userRolesFlag)
				delete(flagValues, dbCollationsFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id without user roles, db collations and db compatibilities",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, userRolesFlag)
				delete(flagValues, dbCollationsFlag)
				delete(flagValues, dbCompatibilitiesFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModelAllTrue(func(model *inputModel) {
				model.UserRoles = false
				model.DBCollations = false
				model.DBCompatibilities = false
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildAndExecuteRequest(t *testing.T) {
	tests := []struct {
		description                string
		model                      *inputModel
		isValid                    bool
		listFlavorsFails           bool
		listVersionsFails          bool
		listStoragesFails          bool
		listUserRolesFails         bool
		listDBCollationsFails      bool
		listDBCompatibilitiesFails bool

		expectListFlavorsCalled           bool
		expectListVersionsCalled          bool
		expectListStoragesCalled          bool
		expectListUserRolesCalled         bool
		expectListDBCollationsCalled      bool
		expectListDBCompatibilitiesCalled bool
	}{
		{
			description:                       "all values",
			model:                             fixtureInputModelAllTrue(),
			isValid:                           true,
			expectListFlavorsCalled:           true,
			expectListVersionsCalled:          true,
			expectListStoragesCalled:          true,
			expectListUserRolesCalled:         true,
			expectListDBCollationsCalled:      true,
			expectListDBCompatibilitiesCalled: true,
		},
		{
			description:                       "no values",
			model:                             fixtureInputModelAllFalse(),
			isValid:                           true,
			expectListFlavorsCalled:           false,
			expectListVersionsCalled:          false,
			expectListStoragesCalled:          false,
			expectListUserRolesCalled:         false,
			expectListDBCollationsCalled:      false,
			expectListDBCompatibilitiesCalled: false,
		},
		{
			description:             "only flavors",
			model:                   fixtureInputModelAllFalse(func(model *inputModel) { model.Flavors = true }),
			isValid:                 true,
			expectListFlavorsCalled: true,
		},
		{
			description:              "only versions",
			model:                    fixtureInputModelAllFalse(func(model *inputModel) { model.Versions = true }),
			isValid:                  true,
			expectListVersionsCalled: true,
		},
		{
			description: "only storages",
			model: fixtureInputModelAllFalse(func(model *inputModel) {
				model.Storages = true
				model.FlavorId = utils.Ptr("2.4")
			}),
			isValid:                  true,
			expectListStoragesCalled: true,
		},
		{
			description: "only user roles",
			model: fixtureInputModelAllFalse(func(model *inputModel) {
				model.UserRoles = true
				model.InstanceId = utils.Ptr(testInstanceId)
			}),
			isValid:                   true,
			expectListUserRolesCalled: true,
		},
		{
			description: "only db collations",
			model: fixtureInputModelAllFalse(func(model *inputModel) {
				model.DBCollations = true
				model.InstanceId = utils.Ptr(testInstanceId)
			}),
			isValid:                      true,
			expectListDBCollationsCalled: true,
		},
		{
			description: "only db compatibilities",
			model: fixtureInputModelAllFalse(func(model *inputModel) {
				model.DBCompatibilities = true
				model.InstanceId = utils.Ptr(testInstanceId)
			}),
			isValid:                           true,
			expectListDBCompatibilitiesCalled: true,
		},
		{
			description:                       "list flavors fails",
			model:                             fixtureInputModelAllTrue(),
			isValid:                           false,
			listFlavorsFails:                  true,
			expectListFlavorsCalled:           true,
			expectListVersionsCalled:          false,
			expectListStoragesCalled:          false,
			expectListUserRolesCalled:         false,
			expectListDBCollationsCalled:      false,
			expectListDBCompatibilitiesCalled: false,
		},
		{
			description:                       "list versions fails",
			model:                             fixtureInputModelAllTrue(),
			isValid:                           false,
			listVersionsFails:                 true,
			expectListFlavorsCalled:           true,
			expectListVersionsCalled:          true,
			expectListStoragesCalled:          false,
			expectListUserRolesCalled:         false,
			expectListDBCollationsCalled:      false,
			expectListDBCompatibilitiesCalled: false,
		},
		{
			description:                       "list storages fails",
			model:                             fixtureInputModelAllTrue(),
			isValid:                           false,
			listStoragesFails:                 true,
			expectListFlavorsCalled:           true,
			expectListVersionsCalled:          true,
			expectListStoragesCalled:          true,
			expectListUserRolesCalled:         false,
			expectListDBCollationsCalled:      false,
			expectListDBCompatibilitiesCalled: false,
		},
		{
			description:                       "list user roles fails",
			model:                             fixtureInputModelAllTrue(),
			isValid:                           false,
			listUserRolesFails:                true,
			expectListFlavorsCalled:           true,
			expectListVersionsCalled:          true,
			expectListStoragesCalled:          true,
			expectListUserRolesCalled:         true,
			expectListDBCollationsCalled:      false,
			expectListDBCompatibilitiesCalled: false,
		},
		{
			description:                       "list db collations fails",
			model:                             fixtureInputModelAllTrue(),
			isValid:                           false,
			listDBCollationsFails:             true,
			expectListFlavorsCalled:           true,
			expectListVersionsCalled:          true,
			expectListStoragesCalled:          true,
			expectListUserRolesCalled:         true,
			expectListDBCollationsCalled:      true,
			expectListDBCompatibilitiesCalled: false,
		},
		{
			description:                       "list db compatibilities fails",
			model:                             fixtureInputModelAllTrue(),
			isValid:                           false,
			listDBCompatibilitiesFails:        true,
			expectListFlavorsCalled:           true,
			expectListVersionsCalled:          true,
			expectListStoragesCalled:          true,
			expectListUserRolesCalled:         true,
			expectListDBCollationsCalled:      true,
			expectListDBCompatibilitiesCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			settings := mockSettings{
				listFlavorsFails:           tt.listFlavorsFails,
				listVersionsFails:          tt.listVersionsFails,
				listStoragesFails:          tt.listStoragesFails,
				listUserRolesFails:         tt.listUserRolesFails,
				listDBCollationsFails:      tt.listDBCollationsFails,
				listDBCompatibilitiesFails: tt.listDBCompatibilitiesFails,
			}

			err := buildAndExecuteRequest(testCtx, params.Printer, tt.model, newAPIMock(&settings))
			if err != nil && tt.isValid {
				t.Fatalf("error building and executing request: %v", err)
			}
			if err == nil && !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if tt.expectListFlavorsCalled != settings.listFlavorsCalled {
				t.Fatalf("expected listFlavorsCalled to be %v, got %v", tt.expectListFlavorsCalled, settings.listFlavorsCalled)
			}
			if tt.expectListVersionsCalled != settings.listVersionsCalled {
				t.Fatalf("expected listVersionsCalled to be %v, got %v", tt.expectListVersionsCalled, settings.listVersionsCalled)
			}
			if tt.expectListStoragesCalled != settings.listStoragesCalled {
				t.Fatalf("expected listStoragesCalled to be %v, got %v", tt.expectListStoragesCalled, settings.listStoragesCalled)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		model             *inputModel
		flavors           *sqlserverflex.ListFlavorsResponse
		versions          *sqlserverflex.ListVersionsResponse
		storages          *sqlserverflex.ListStoragesResponse
		userRoles         *sqlserverflex.ListRolesResponse
		dbCollations      *sqlserverflex.ListCollationsResponse
		dbCompatibilities *sqlserverflex.ListCompatibilityResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty - only model",
			args: args{
				model: fixtureInputModelAllFalse(),
			},
			wantErr: false,
		},
		{
			name: "all input set",
			args: args{
				model:             fixtureInputModelAllTrue(),
				flavors:           &sqlserverflex.ListFlavorsResponse{Flavors: []sqlserverflex.ListFlavors{}},
				versions:          &sqlserverflex.ListVersionsResponse{Versions: []sqlserverflex.Version{}},
				storages:          &sqlserverflex.ListStoragesResponse{StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{}},
				userRoles:         &sqlserverflex.ListRolesResponse{Roles: []string{}},
				dbCollations:      &sqlserverflex.ListCollationsResponse{Collations: []sqlserverflex.DatabaseGetcollation{}},
				dbCompatibilities: &sqlserverflex.ListCompatibilityResponse{Compatibilities: []sqlserverflex.DatabaseGetcompatibility{}},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.flavors, tt.args.versions, tt.args.storages, tt.args.userRoles, tt.args.dbCollations, tt.args.dbCompatibilities); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
