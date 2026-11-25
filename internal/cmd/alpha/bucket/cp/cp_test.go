package cp

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	testRegion   = "eu01"
	testSrc      = "./local/file.txt"
	testDst      = "s3://my-bucket/file.txt"
	testSrcRef   = "s3://my-bucket/file.txt"
	testDstLocal = "./local/file.txt"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		Source:      testSrc,
		Destination: testDst,
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
			description:   "upload",
			argValues:     []string{testSrc, testDst},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "download",
			argValues:   []string{testSrcRef, testDstLocal},
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Source = testSrcRef
				model.Destination = testDstLocal
			}),
		},
		{
			description: "with recursive",
			argValues:   []string{testSrc, testDst},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[recursiveFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Recursive = true
			}),
		},
		{
			description: "too few args",
			argValues:   []string{testSrc},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "too many args",
			argValues:   []string{testSrc, testDst, "extra"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestOutputResult(t *testing.T) {
	params := testparams.NewTestParams()
	if err := outputResult(params.Printer, "", "Uploaded 1 object(s) to \"s3://my-bucket/file.txt\"", 1); err != nil {
		t.Fatalf("outputResult() unexpected error: %v", err)
	}
	if output := params.Out.String(); output == "" {
		t.Fatalf("outputResult() expected non-empty output")
	}
}
