package delete

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const testProfile = "test-profile"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testProfile,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		Profile: testProfile,
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
			isValid:     false,
		},
		{
			description: "some global flag",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.VerbosityFlag: globalflags.DebugVerbosity,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Verbosity = globalflags.DebugVerbosity
			}),
		},
		{
			description: "invalid profile",
			argValues:   []string{"invalid-profile-&"},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestDeleteProfileWithoutKeyring(t *testing.T) {
	params := testparams.NewTestParams()
	params.Printer.AssumeYes = true
	profile := fmt.Sprintf("test-profile-%s", time.Now().Format("20060102150405"))
	path := config.GetProfileFolderPath(profile)
	t.Cleanup(func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("cleanup: remove profile folder at path %q: %v", path, err)
		}
	})
	err := config.ValidateProfile(profile)
	if err != nil {
		t.Fatalf("validate profile %q: %v", profile, err)
	}
	err = config.CreateProfile(params.Printer, profile, true, false, true)
	if err != nil {
		t.Fatalf("create profile %q: %v", profile, err)
	}
	keyring.MockInitWithError(keyring.ErrUnsupportedPlatform)
	deleteCmd := NewCmd(params.CmdParams)
	err = deleteCmd.RunE(deleteCmd, []string{profile})
	if err != nil {
		t.Fatalf("run cmd: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected profile folder to be deleted, but it still exists at path %q", path)
	}
}
