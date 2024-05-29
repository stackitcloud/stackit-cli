package config

import (
	"path/filepath"
	"testing"
)

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
