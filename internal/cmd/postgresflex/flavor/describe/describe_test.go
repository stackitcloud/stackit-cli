package describe

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

var (
	testProjectId = uuid.NewString()
)

const (
	testRegion   = "eu01"
	testFlavorId = "16.64"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testFlavorId,
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
		FlavorId: testFlavorId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureTestFlavor() postgresflex.ListFlavors {
	return postgresflex.ListFlavors{
		Id:          testFlavorId,
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
	}
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
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
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
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

type mockSettings struct {
	listFlavorsFails bool
	flavors          []postgresflex.ListFlavors
}

func newApiMock(s mockSettings) postgresflex.DefaultAPI {
	return postgresflex.DefaultAPIServiceMock{
		ListFlavorsExecuteMock: utils.Ptr(func(_ postgresflex.ApiListFlavorsRequest) (*postgresflex.ListFlavorsResponse, error) {
			if s.listFlavorsFails {
				return nil, fmt.Errorf("mock list flavors fails")
			}

			return &postgresflex.ListFlavorsResponse{
				Flavors: s.flavors,
			}, nil
		}),
	}
}

func TestBuildAndExecuteRequest(t *testing.T) {
	tests := []struct {
		description  string
		model        *inputModel
		mockSettings mockSettings
		want         *postgresflex.ListFlavors
		wantErr      bool
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			mockSettings: mockSettings{
				flavors: []postgresflex.ListFlavors{
					fixtureTestFlavor(),
				},
			},
			want:    utils.Ptr(fixtureTestFlavor()),
			wantErr: false,
		},
		{
			description: "flavor not found",
			model: fixtureInputModel(func(model *inputModel) {
				model.FlavorId = "foo-bar-does-not-exist"
			}),
			mockSettings: mockSettings{
				flavors: []postgresflex.ListFlavors{
					fixtureTestFlavor(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			description: "list flavors fails",
			model:       fixtureInputModel(),
			mockSettings: mockSettings{
				listFlavorsFails: true,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual, err := buildAndExecuteRequest(context.Background(), tt.model, newApiMock(tt.mockSettings))
			if err != nil && !tt.wantErr {
				t.Errorf("want no error, got %v", err)
			}

			if err == nil && tt.wantErr {
				t.Errorf("want error, got nil")
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
		outputFormat string
		flavor       *postgresflex.ListFlavors
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "flavor is nil",
			args: args{
				flavor: nil,
			},
			wantErr: true,
		},
		{
			name: "flavor is empty",
			args: args{
				flavor: &postgresflex.ListFlavors{},
			},
			wantErr: false,
		},
		{
			name: "complete",
			args: args{
				flavor: utils.Ptr(fixtureTestFlavor()),
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.flavor); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
