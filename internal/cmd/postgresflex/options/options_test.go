package options

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testProjectId = uuid.NewString()

type mockSettings struct {
	listFlavorsFails  bool
	listVersionsFails bool
	listStoragesFails bool

	listFlavorsCalled  bool
	listVersionsCalled bool
	listStoragesCalled bool
}

func newAPIClientMock(c *mockSettings) postgresflex.DefaultAPI {
	return postgresflex.DefaultAPIServiceMock{
		ListFlavorsExecuteMock: utils.Ptr(func(_ postgresflex.ApiListFlavorsRequest) (*postgresflex.ListFlavorsResponse, error) {
			c.listFlavorsCalled = true
			if c.listFlavorsFails {
				return nil, fmt.Errorf("list flavors failed")
			}
			return utils.Ptr(postgresflex.ListFlavorsResponse{
				Flavors: []postgresflex.Flavor{},
			}), nil
		}),
		ListVersionsExecuteMock: utils.Ptr(func(_ postgresflex.ApiListVersionsRequest) (*postgresflex.ListVersionsResponse, error) {
			c.listVersionsCalled = true
			if c.listVersionsFails {
				return nil, fmt.Errorf("list versions failed")
			}
			return utils.Ptr(postgresflex.ListVersionsResponse{
				Versions: []string{},
			}), nil
		}),
		ListStoragesExecuteMock: utils.Ptr(func(_ postgresflex.ApiListStoragesRequest) (*postgresflex.ListStoragesResponse, error) {
			c.listStoragesCalled = true
			if c.listStoragesFails {
				return nil, fmt.Errorf("list storages failed")
			}
			return utils.Ptr(postgresflex.ListStoragesResponse{
				StorageClasses: []string{},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			}), nil
		}),
	}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		flavorsFlag:               "true",
		versionsFlag:              "true",
		storagesFlag:              "true",
		flavorIdFlag:              "2.4",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModelAllFalse(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		Flavors:  false,
		Versions: false,
		Storages: false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureInputModelAllTrue(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		Flavors:  true,
		Versions: true,
		Storages: true,
		FlavorId: utils.Ptr("2.4"),
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
			expectedModel: fixtureInputModelAllFalse(func(model *inputModel) {
				model.Flavors = true
				model.Versions = true
			}),
		},
		{
			description: "some values 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, flavorsFlag)
				delete(flagValues, versionsFlag)
				flagValues[storagesFlag] = "true"
				flagValues[flavorIdFlag] = "2.4"
			}),
			isValid: true,
			expectedModel: fixtureInputModelAllFalse(func(model *inputModel) {
				model.Storages = true
				model.FlavorId = utils.Ptr("2.4")
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildAndExecuteRequest(t *testing.T) {
	tests := []struct {
		description              string
		model                    *inputModel
		isValid                  bool
		mockClientSettings       mockSettings
		expectListFlavorsCalled  bool
		expectListVersionsCalled bool
		expectListStoragesCalled bool
	}{
		{
			description:              "all values",
			model:                    fixtureInputModelAllTrue(),
			isValid:                  true,
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: true,
			expectListStoragesCalled: true,
		},
		{
			description:              "no values",
			model:                    fixtureInputModelAllFalse(),
			isValid:                  true,
			expectListFlavorsCalled:  false,
			expectListVersionsCalled: false,
			expectListStoragesCalled: false,
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
			description: "list flavors fails",
			model:       fixtureInputModelAllTrue(),
			isValid:     false,
			mockClientSettings: mockSettings{
				listFlavorsFails: true,
			},
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: false,
			expectListStoragesCalled: false,
		},
		{
			description: "list versions fails",
			model:       fixtureInputModelAllTrue(),
			isValid:     false,
			mockClientSettings: mockSettings{
				listVersionsFails: true,
			},
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: true,
			expectListStoragesCalled: false,
		},
		{
			description: "list storages fails",
			model:       fixtureInputModelAllTrue(),
			isValid:     false,
			mockClientSettings: mockSettings{
				listStoragesFails: true,
			},
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: true,
			expectListStoragesCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			client := newAPIClientMock(&tt.mockClientSettings)
			err := buildAndExecuteRequest(testCtx, params.Printer, tt.model, client)
			if err != nil && tt.isValid {
				t.Fatalf("error building and executing request: %v", err)
			}
			if err == nil && !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if tt.expectListFlavorsCalled != tt.mockClientSettings.listFlavorsCalled {
				t.Fatalf("expected listFlavorsCalled to be %v, got %v", tt.expectListFlavorsCalled, tt.mockClientSettings.listFlavorsCalled)
			}
			if tt.expectListVersionsCalled != tt.mockClientSettings.listVersionsCalled {
				t.Fatalf("expected listVersionsCalled to be %v, got %v", tt.expectListVersionsCalled, tt.mockClientSettings.listVersionsCalled)
			}
			if tt.expectListStoragesCalled != tt.mockClientSettings.listStoragesCalled {
				t.Fatalf("expected listStoragesCalled to be %v, got %v", tt.expectListStoragesCalled, tt.mockClientSettings.listStoragesCalled)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		model    inputModel
		flavors  *postgresflex.ListFlavorsResponse
		versions *postgresflex.ListVersionsResponse
		storages *postgresflex.ListStoragesResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{model: inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}}}, false},
		{"standard", args{
			model:    inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			flavors:  &postgresflex.ListFlavorsResponse{},
			versions: &postgresflex.ListVersionsResponse{},
			storages: &postgresflex.ListStoragesResponse{},
		}, false},
		{
			"complete",
			args{
				model: inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}, Flavors: false, Versions: false, Storages: false, FlavorId: new(string)},
				flavors: &postgresflex.ListFlavorsResponse{
					Flavors: []postgresflex.Flavor{},
				},
				versions: &postgresflex.ListVersionsResponse{
					Versions: []string{},
				},
				storages: &postgresflex.ListStoragesResponse{
					StorageClasses: []string{},
					StorageRange:   &postgresflex.StorageRange{},
				},
			},
			false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.flavors, tt.args.versions, tt.args.storages); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
