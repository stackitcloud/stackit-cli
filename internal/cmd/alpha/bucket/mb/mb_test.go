package mb

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	testRegion = "eu01"
	testBucket = "s3://my-bucket"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{testBucket}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

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
		URI: testBucket,
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid uri, wrong scheme",
			argValues:   []string{"https://my-bucket"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid uri, object instead of bucket",
			argValues:   []string{"s3://my-bucket/key"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid uri, invalid bucket name",
			argValues:   []string{"s3://My-Bucket"},
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
	if err := outputResult(params.Printer, "", testBucket); err != nil {
		t.Fatalf("outputResult() unexpected error: %v", err)
	}
	if output := params.Out.String(); output == "" {
		t.Fatalf("outputResult() expected non-empty output")
	}
}

func TestLocationConstraint(t *testing.T) {
	tests := []struct {
		description string
		region      string
		expected    string
	}{
		{description: "empty region omits the constraint", region: "", expected: ""},
		{description: "default region omits the constraint", region: defaultRegion, expected: ""},
		{description: "other region", region: "eu02", expected: "eu02"},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if got := locationConstraint(tt.region); got != tt.expected {
				t.Fatalf("locationConstraint(%q) = %q, want %q", tt.region, got, tt.expected)
			}
		})
	}
}
