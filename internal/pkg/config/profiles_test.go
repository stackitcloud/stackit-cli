package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

//go:embed template/test_profile.json
var templateConfig string

func TestValidateProfile(t *testing.T) {
	tests := []struct {
		description string
		profile     string
		isValid     bool
	}{
		{
			description: "valid profile with letters",
			profile:     "myprofile",
			isValid:     true,
		},
		{
			description: "valid with letters and hyphen",
			profile:     "my-profile",
			isValid:     true,
		},
		{
			description: "valid with letters, numbers, and hyphen",
			profile:     "my-profile-123",
			isValid:     true,
		},
		{
			description: "valid with letters, numbers, and ending with hyphen",
			profile:     "my-profile123-",
			isValid:     true,
		},
		{
			description: "invalid empty",
			profile:     "",
			isValid:     false,
		},
		{
			description: "invalid with special characters",
			profile:     "my_profile",
			isValid:     false,
		},
		{
			description: "invalid with spaces",
			profile:     "my profile",
			isValid:     false,
		},
		{
			description: "invalid starting with -",
			profile:     "-my-profile",
			isValid:     false,
		},
		{
			description: "invalid profile with uppercase letters",
			profile:     "myProfile",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateProfile(tt.profile)
			if tt.isValid && err != nil {
				t.Errorf("expected profile to be valid but got error: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("expected profile to be invalid but got no error")
			}
		})
	}
}

func TestGetProfileFolderPath(t *testing.T) {
	tests := []struct {
		description               string
		defaultConfigFolderNotSet bool
		profile                   string
		expected                  string
	}{
		{
			description: "default profile",
			profile:     DefaultProfileName,
			expected:    getInitialConfigDir(),
		},
		{
			description:               "default profile, default config folder not set",
			defaultConfigFolderNotSet: true,
			profile:                   DefaultProfileName,
			expected:                  getInitialConfigDir(),
		},
		{
			description: "custom profile",
			profile:     "my-profile",
			expected:    filepath.Join(getInitialConfigDir(), profileRootFolder, "my-profile"),
		},
		{
			description:               "custom profile, default config folder not set",
			defaultConfigFolderNotSet: true,
			profile:                   "my-profile",
			expected:                  filepath.Join(getInitialConfigDir(), profileRootFolder, "my-profile"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			defaultConfigFolderPath = getInitialConfigDir()
			if tt.defaultConfigFolderNotSet {
				defaultConfigFolderPath = ""
			}
			actual := GetProfileFolderPath(tt.profile)
			if actual != tt.expected {
				t.Errorf("expected profile folder path to be %q but got %q", tt.expected, actual)
			}
		})
	}
}

func TestImportProfile(t *testing.T) {
	tests := []struct {
		description string
		profile     string
		config      string
		setAsActive bool
		isValid     bool
	}{
		{
			description: "valid profile",
			profile:     "profile-name",
			config:      templateConfig,
			setAsActive: false,
			isValid:     true,
		},
		{
			description: "invalid profile name",
			profile:     "invalid-profile-&",
			config:      templateConfig,
			setAsActive: false,
			isValid:     false,
		},
		{
			description: "invalid config",
			profile:     "my-profile",
			config:      `{ "invalid": "json }`,
			setAsActive: false,
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			err := ImportProfile(p, tt.profile, tt.config, tt.setAsActive)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("profile should be valid but got error: %v\n", err)
			}

			if !tt.isValid {
				t.Fatalf("profile should be invalid but got no error\n")
			}
		})

		t.Cleanup(func() {
			p := print.NewPrinter()
			err := DeleteProfile(p, tt.profile)
			if err != nil {
				if !tt.isValid {
					return
				}
				fmt.Printf("could not clean up imported profile: %v\n", err)
			}
		})
	}
}

func TestExportProfile(t *testing.T) {
	// Create directory where the export configs should be stored
	testDir, err := os.MkdirTemp(os.TempDir(), "stackit-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Printf("could not clean up temp dir: %v\n", err)
			}
		}(testDir)
	})

	defaultConfigFolderPath = filepath.Join(testDir, "config")
	err = os.Mkdir(defaultConfigFolderPath, 0o750)
	if err != nil {
		t.Fatal(err)
	}

	// Create prerequisite profile
	p := print.NewPrinter()
	profileName := "export-profile-test"
	err = CreateProfile(p, profileName, true, true)
	if err != nil {
		t.Fatalf("could not create prerequisite profile, %v", err)
	}
	initConfig(defaultConfigFolderPath)
	err = Write()
	if err != nil {
		t.Fatalf("could not write profile, %v", err)
	}
	t.Cleanup(func() {
		func(p *print.Printer, profile string) {
			err := DeleteProfile(p, profile)
			if err != nil {
				fmt.Printf("could not clean up prerequisite profile %q, %v", profileName, err)
			}
		}(p, profileName)
	})

	tests := []struct {
		description string
		profile     string
		filePath    string
		isValid     bool
	}{
		{
			description: "valid profile",
			profile:     profileName,
			filePath:    filepath.Join(testDir, fmt.Sprintf("custom-name.%s", configFileExtension)),
			isValid:     true,
		},
		{
			description: "invalid profile",
			profile:     "invalid-my-profile",
			isValid:     false,
		},
		{
			description: "not existing path",
			profile:     profileName,
			filePath:    filepath.Join(testDir, "invalid", "path", fmt.Sprintf("custom-name.%s", configFileExtension)),
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			err := ExportProfile(p, tt.profile, tt.filePath)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("export should be valid but got error: %v\n", err)
			}
			if !tt.isValid {
				t.Fatalf("export should be invalid but got no error\n")
			}
		})
	}
}
