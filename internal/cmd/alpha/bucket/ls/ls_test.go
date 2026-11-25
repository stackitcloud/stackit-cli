package ls

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/bucket"
)

const (
	testRegion = "eu01"
	testURI    = "s3://my-bucket/prefix/"
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
			description:   "no uri",
			argValues:     []string{},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "with uri",
			argValues:   []string{testURI},
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.URI = testURI
			}),
		},
		{
			description: "with recursive",
			argValues:   []string{testURI},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[recursiveFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.URI = testURI
				model.Recursive = true
			}),
		},
		{
			description: "too many args",
			argValues:   []string{testURI, "extra"},
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

func TestOutputResultObjects(t *testing.T) {
	params := testparams.NewTestParams()
	prefixes := []string{"dir1/"}
	objects := []bucket.ObjectInfo{
		{Key: "dir1/file.txt", Size: 123},
	}
	if err := outputResultObjects(params.Printer, "", prefixes, objects); err != nil {
		t.Fatalf("outputResultObjects() unexpected error: %v", err)
	}
	if output := params.Out.String(); output == "" {
		t.Fatalf("outputResultObjects() expected non-empty output")
	}
}

func TestOutputResultBuckets(t *testing.T) {
	params := testparams.NewTestParams()
	buckets := []bucket.BucketInfo{
		{Name: "my-bucket"},
	}
	if err := outputResultBuckets(params.Printer, "", buckets); err != nil {
		t.Fatalf("outputResultBuckets() unexpected error: %v", err)
	}
	if output := params.Out.String(); output == "" {
		t.Fatalf("outputResultBuckets() expected non-empty output")
	}
}
