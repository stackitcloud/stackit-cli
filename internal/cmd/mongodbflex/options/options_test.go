package options

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")

type mongoDBFlexClientMocked struct {
	listFlavorsFails  bool
	listVersionsFails bool
	listStoragesFails bool

	listFlavorsCalled  bool
	listVersionsCalled bool
	listStoragesCalled bool
}

func (c *mongoDBFlexClientMocked) ListFlavorsExecute(_ context.Context, _, _ string) (*mongodbflex.ListFlavorsResponse, error) {
	c.listFlavorsCalled = true
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list flavors failed")
	}
	return utils.Ptr(mongodbflex.ListFlavorsResponse{
		Flavors: utils.Ptr([]mongodbflex.InstanceFlavor{}),
	}), nil
}

func (c *mongoDBFlexClientMocked) ListVersionsExecute(_ context.Context, _, _ string) (*mongodbflex.ListVersionsResponse, error) {
	c.listVersionsCalled = true
	if c.listVersionsFails {
		return nil, fmt.Errorf("list versions failed")
	}
	return utils.Ptr(mongodbflex.ListVersionsResponse{
		Versions: utils.Ptr([]string{}),
	}), nil
}

func (c *mongoDBFlexClientMocked) ListStoragesExecute(_ context.Context, _, _, _ string) (*mongodbflex.ListStoragesResponse, error) {
	c.listStoragesCalled = true
	if c.listStoragesFails {
		return nil, fmt.Errorf("list storages failed")
	}
	return utils.Ptr(mongodbflex.ListStoragesResponse{
		StorageClasses: utils.Ptr([]string{}),
		StorageRange: &mongodbflex.StorageRange{
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
			cmd := NewCmd(&types.CmdParams{Printer: p})
			p.Cmd = cmd
			client := &mongoDBFlexClientMocked{
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

func TestOutputResult(t *testing.T) {
	type args struct {
		inputModel *inputModel
		flavors    *mongodbflex.ListFlavorsResponse
		versions   *mongodbflex.ListVersionsResponse
		storages   *mongodbflex.ListStoragesResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "missing model",
			args: args{
				flavors:  &mongodbflex.ListFlavorsResponse{},
				versions: &mongodbflex.ListVersionsResponse{},
				storages: &mongodbflex.ListStoragesResponse{},
			},
			wantErr: true,
		},
		{
			name: "empty model",
			args: args{
				inputModel: &inputModel{},
				flavors:    &mongodbflex.ListFlavorsResponse{},
				versions:   &mongodbflex.ListVersionsResponse{},
				storages:   &mongodbflex.ListStoragesResponse{},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				inputModel: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				flavors:  &mongodbflex.ListFlavorsResponse{},
				versions: &mongodbflex.ListVersionsResponse{},
				storages: &mongodbflex.ListStoragesResponse{},
			},
			wantErr: false,
		},
		{
			name: "missing flavors",
			args: args{
				inputModel: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				versions: &mongodbflex.ListVersionsResponse{},
				storages: &mongodbflex.ListStoragesResponse{},
			},
			wantErr: false,
		},
		{
			name: "missing versions",
			args: args{
				inputModel: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				flavors:  &mongodbflex.ListFlavorsResponse{},
				storages: &mongodbflex.ListStoragesResponse{},
			},
			wantErr: false,
		},
		{
			name: "missing storages",
			args: args{
				inputModel: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{},
				},
				flavors:  &mongodbflex.ListFlavorsResponse{},
				versions: &mongodbflex.ListVersionsResponse{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.inputModel, tt.args.flavors, tt.args.versions, tt.args.storages); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutputResultAsTable(t *testing.T) {
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
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "missing input model",
			args: args{
				options: &options{},
			},
			wantErr: true,
		},
		{
			name: "missing options",
			args: args{
				model: &inputModel{},
			},
			wantErr: true,
		},
		{
			name: "empty input model and empty options",
			args: args{
				model:   &inputModel{},
				options: &options{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResultAsTable(p, tt.args.model, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("outputResultAsTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
