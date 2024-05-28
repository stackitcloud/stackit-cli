package options

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")

type sqlServerFlexClientMocked struct {
	listFlavorsFails  bool
	listVersionsFails bool
	listStoragesFails bool

	listFlavorsCalled  bool
	listVersionsCalled bool
	listStoragesCalled bool
}

func (c *sqlServerFlexClientMocked) ListFlavorsExecute(_ context.Context, _ string) (*sqlserverflex.ListFlavorsResponse, error) {
	c.listFlavorsCalled = true
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list flavors failed")
	}
	return utils.Ptr(sqlserverflex.ListFlavorsResponse{
		Flavors: utils.Ptr([]sqlserverflex.InstanceFlavorEntry{}),
	}), nil
}

func (c *sqlServerFlexClientMocked) ListVersionsExecute(_ context.Context, _ string) (*sqlserverflex.ListVersionsResponse, error) {
	c.listVersionsCalled = true
	if c.listVersionsFails {
		return nil, fmt.Errorf("list versions failed")
	}
	return utils.Ptr(sqlserverflex.ListVersionsResponse{
		Versions: utils.Ptr([]string{}),
	}), nil
}

func (c *sqlServerFlexClientMocked) ListStoragesExecute(_ context.Context, _, _ string) (*sqlserverflex.ListStoragesResponse, error) {
	c.listStoragesCalled = true
	if c.listStoragesFails {
		return nil, fmt.Errorf("list storages failed")
	}
	return utils.Ptr(sqlserverflex.ListStoragesResponse{
		StorageClasses: utils.Ptr([]string{}),
		StorageRange: &sqlserverflex.StorageRange{
			Min: utils.Ptr(int64(10)),
			Max: utils.Ptr(int64(100)),
		},
	}), nil
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		flavorsFlag:  "true",
		versionsFlag: "true",
		storagesFlag: "true",
		flavorIdFlag: "2.4",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModelAllFalse(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		Flavors:         false,
		Versions:        false,
		Storages:        false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureInputModelAllTrue(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		Flavors:         true,
		Versions:        true,
		Storages:        true,
		FlavorId:        utils.Ptr("2.4"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
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
			p := print.NewPrinter()
			cmd := NewCmd(p)
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

func TestBuildAndExecuteRequest(t *testing.T) {
	tests := []struct {
		description              string
		model                    *inputModel
		isValid                  bool
		listFlavorsFails         bool
		listVersionsFails        bool
		listStoragesFails        bool
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
			description:              "list flavors fails",
			model:                    fixtureInputModelAllTrue(),
			isValid:                  false,
			listFlavorsFails:         true,
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: false,
			expectListStoragesCalled: false,
		},
		{
			description:              "list versions fails",
			model:                    fixtureInputModelAllTrue(),
			isValid:                  false,
			listVersionsFails:        true,
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: true,
			expectListStoragesCalled: false,
		},
		{
			description:              "list storages fails",
			model:                    fixtureInputModelAllTrue(),
			isValid:                  false,
			listStoragesFails:        true,
			expectListFlavorsCalled:  true,
			expectListVersionsCalled: true,
			expectListStoragesCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := &print.Printer{}
			cmd := NewCmd(p)
			p.Cmd = cmd
			client := &sqlServerFlexClientMocked{
				listFlavorsFails:  tt.listFlavorsFails,
				listVersionsFails: tt.listVersionsFails,
				listStoragesFails: tt.listStoragesFails,
			}

			err := buildAndExecuteRequest(testCtx, p, tt.model, client)
			if err != nil && tt.isValid {
				t.Fatalf("error building and executing request: %v", err)
			}
			if err == nil && !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if tt.expectListFlavorsCalled != client.listFlavorsCalled {
				t.Fatalf("expected listFlavorsCalled to be %v, got %v", tt.expectListFlavorsCalled, client.listFlavorsCalled)
			}
			if tt.expectListVersionsCalled != client.listVersionsCalled {
				t.Fatalf("expected listVersionsCalled to be %v, got %v", tt.expectListVersionsCalled, client.listVersionsCalled)
			}
			if tt.expectListStoragesCalled != client.listStoragesCalled {
				t.Fatalf("expected listStoragesCalled to be %v, got %v", tt.expectListStoragesCalled, client.listStoragesCalled)
			}
		})
	}
}
