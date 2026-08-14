package options

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

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

	flavors  []postgresflex.ListFlavors
	versions []postgresflex.Version

	listFlavorsCalled  bool
	listVersionsCalled bool
}

func newAPIClientMock(c *mockSettings) postgresflex.DefaultAPI {
	return postgresflex.DefaultAPIServiceMock{
		ListFlavorsExecuteMock: utils.Ptr(func(_ postgresflex.ApiListFlavorsRequest) (*postgresflex.ListFlavorsResponse, error) {
			c.listFlavorsCalled = true

			if c.listFlavorsFails {
				return nil, fmt.Errorf("list flavors failed")
			}

			return &postgresflex.ListFlavorsResponse{
				Flavors: c.flavors,
			}, nil
		}),
		ListVersionsExecuteMock: utils.Ptr(func(_ postgresflex.ApiListVersionsRequest) (*postgresflex.ListVersionsResponse, error) {
			c.listVersionsCalled = true

			if c.listVersionsFails {
				return nil, fmt.Errorf("list versions failed")
			}

			return &postgresflex.ListVersionsResponse{
				Versions: c.versions,
			}, nil
		}),
	}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
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
		Flavors:  false,
		Versions: false,
		Storages: false,
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
			description: "default (invalid) - at least one flag must be set",
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "all values",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[versionsFlag] = "true"
				flagValues[storagesFlag] = "true"
				flagValues[flavorsFlag] = "true"
				flagValues[flavorIdFlag] = "16.64"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Storages = true
				model.Versions = true
				model.Flavors = true
				model.FlavorId = utils.Ptr("16.64")
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "some values 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[versionsFlag] = "true"
				flagValues[flavorsFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
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
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Storages = true
				model.FlavorId = utils.Ptr("2.4")
			}),
		},
		{
			description: "storages without flavor-id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[storagesFlag] = "true"
				delete(flagValues, flavorIdFlag)
			}),
			isValid: false,
		},
		{
			description: "flavor-id without storage",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[storagesFlag] = "false"
				delete(flagValues, storagesFlag)
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

func TestBuildAndExecuteRequest(t *testing.T) {
	fixtureVersions := func() []postgresflex.Version {
		return []postgresflex.Version{
			{
				Version:    "2.4",
				Beta:       false,
				Recommend:  true,
				Deprecated: "",
			},
			{
				Version:    "1.9",
				Beta:       false,
				Recommend:  false,
				Deprecated: time.Date(2026, time.July, 27, 12, 0, 0, 0, time.UTC).String(),
			},
			{
				Version:   "3.0",
				Beta:      true,
				Recommend: false,
			},
		}
	}

	fixtureFlavors := func() []postgresflex.ListFlavors {
		return []postgresflex.ListFlavors{
			{
				Id:          "16.64",
				Cpu:         16,
				Description: "PostgreSQL-Flex-16.64-Single-EU01",
				MaxGB:       4000,
				MinGB:       5,
				NodeType:    "Single",
				StorageClasses: []postgresflex.FlavorStorageClassesStorageClass{
					{
						Class:          "premium-perf2-stackit",
						MaxIoPerSec:    1000,
						MaxThroughInMb: 100,
					},
					{
						Class:          "premium-perf4-stackit",
						MaxIoPerSec:    2000,
						MaxThroughInMb: 150,
					},
					{
						Class:          "premium-perf6-stackit",
						MaxIoPerSec:    5000,
						MaxThroughInMb: 200,
					},
				},
			},
		}
	}

	tests := []struct {
		description              string
		model                    *inputModel
		isValid                  bool
		mockClientSettings       mockSettings
		expectListFlavorsCalled  bool
		expectListVersionsCalled bool
		want                     *options
	}{
		{
			description: "all values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Storages = true
				model.Versions = true
				model.Flavors = true
			}),
			mockClientSettings: mockSettings{
				flavors:  fixtureFlavors(),
				versions: fixtureVersions(),
			},
			isValid:                  true,
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: true,
			want: &options{
				Flavors:  fixtureFlavors(),
				Versions: fixtureVersions(),
			},
		},
		{
			description:              "no values",
			model:                    fixtureInputModel(),
			isValid:                  true,
			expectListFlavorsCalled:  false,
			expectListVersionsCalled: false,
			want:                     &options{},
		},
		{
			description:             "only flavors",
			model:                   fixtureInputModel(func(model *inputModel) { model.Flavors = true }),
			isValid:                 true,
			expectListFlavorsCalled: true,
			want:                    &options{},
		},
		{
			description: "only versions",
			model:       fixtureInputModel(func(model *inputModel) { model.Versions = true }),
			mockClientSettings: mockSettings{
				versions: fixtureVersions(),
			},
			isValid:                  true,
			expectListVersionsCalled: true,
			want: &options{
				Versions: fixtureVersions(),
			},
		},
		{
			description: "only storages - flavor not found",
			model: fixtureInputModel(func(model *inputModel) {
				model.Storages = true
				model.FlavorId = utils.Ptr("2.2")
			}),
			isValid: false,
		},
		{
			description: "only storages - flavor found",
			model: fixtureInputModel(func(model *inputModel) {
				model.Storages = true
				model.FlavorId = func() *string {
					return &fixtureFlavors()[0].Id
				}()
			}),
			mockClientSettings: mockSettings{
				flavors: fixtureFlavors(),
			},
			expectListFlavorsCalled: true,
			want: &options{
				Storages: func() *flavorStorages {
					return &flavorStorages{
						FlavorId: fixtureFlavors()[0].Id,
						Storages: fixtureFlavors()[0].StorageClasses,
					}
				}(),
			},
			isValid: true,
		},
		{
			description: "list flavors fails",
			model: fixtureInputModel(func(model *inputModel) {
				model.Storages = true
				model.Versions = true
				model.Flavors = true
			}),
			isValid: false,
			mockClientSettings: mockSettings{
				listFlavorsFails: true,
			},
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: false,
		},
		{
			description: "list versions fails",
			model: fixtureInputModel(func(model *inputModel) {
				model.Storages = true
				model.Versions = true
				model.Flavors = true
			}),
			isValid: false,
			mockClientSettings: mockSettings{
				listVersionsFails: true,
			},
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual, err := buildAndExecuteRequest(testCtx, tt.model, newAPIClientMock(&tt.mockClientSettings))
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

			diff := cmp.Diff(actual, tt.want)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		model   *inputModel
		options *options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "options is nil",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				options: nil,
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				options: &options{},
			},
			wantErr: false,
		},
		{
			name: "empty flavors and versions slice",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				options: &options{
					Flavors:  []postgresflex.ListFlavors{},
					Versions: []postgresflex.Version{},
				},
			},
			wantErr: false,
		},
		{
			name: "complete",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
					Flavors:         false,
					Versions:        false,
					Storages:        false,
					FlavorId:        new(string),
				},
				options: &options{
					Flavors:  []postgresflex.ListFlavors{},
					Versions: []postgresflex.Version{},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
