package mv

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/bucket"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	testRegion = "eu01"
	testSrc    = "./local/file.txt"
	testDst    = "s3://my-bucket/file.txt"
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
			description:   "base",
			argValues:     []string{testSrc, testDst},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestOutputResult(t *testing.T) {
	params := testparams.NewTestParams()
	if err := outputResult(params.Printer, "", "Moved 1 object(s) from \"./local/file.txt\" to \"s3://my-bucket/file.txt\"", 1); err != nil {
		t.Fatalf("outputResult() unexpected error: %v", err)
	}
	if output := params.Out.String(); output == "" {
		t.Fatalf("outputResult() expected non-empty output")
	}
}

func TestIsCrossEndpointMove(t *testing.T) {
	srcEU01 := bucket.Reference{Bucket: "src", Endpoint: "https://object.storage.eu01.onstackit.cloud", Region: "eu01"}
	dstEU01 := bucket.Reference{Bucket: "dst", Endpoint: "https://object.storage.eu01.onstackit.cloud", Region: "eu01"}
	dstEU02 := bucket.Reference{Bucket: "dst", Endpoint: "https://object.storage.eu02.onstackit.cloud", Region: "eu02"}
	bare := bucket.Reference{Bucket: "dst"}

	tests := []struct {
		description string
		src         bucket.Reference
		dst         bucket.Reference
		expected    bool
	}{
		{description: "same region", src: srcEU01, dst: dstEU01, expected: false},
		{description: "different regions", src: srcEU01, dst: dstEU02, expected: true},
		{description: "different endpoints, same region", src: bucket.Reference{Endpoint: "https://a.example"}, dst: bucket.Reference{Endpoint: "https://b.example"}, expected: true},
		{description: "one side without region", src: srcEU01, dst: bare, expected: false},
		{description: "both plain references", src: bare, dst: bare, expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if got := isCrossEndpointMove(tt.src, tt.dst); got != tt.expected {
				t.Fatalf("isCrossEndpointMove = %v, want %v", got, tt.expected)
			}
		})
	}
}
