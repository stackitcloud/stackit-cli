package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// GetProfile returns the current profile to be used by the CLI.
//
// The profile is determined by the value of the STACKIT_CLI_PROFILE environment variable, or, if not set,
// by the contents of the profile file in the CLI config folder.
//
// If the environment variable is not set and the profile file does not exist, it returns an empty string.
//
// If the profile is not valid, it returns an error.
func GetProfile() (string, error) {
	profile, profileSet := os.LookupEnv("STACKIT_CLI_PROFILE")
	if !profileSet {
		profileFilePath := filepath.Join(configFolderPath, fmt.Sprintf("%s.%s", profileFileName, profileFileExtension))
		contents, exists, err := readFileIfExists(profileFilePath)
		if err != nil {
			return "", fmt.Errorf("read profile from file: %v", err)
		}
		if !exists {
			return "", nil
		}
		profile = contents
	}

	err := validateProfile(profile)
	if err != nil {
		return "", fmt.Errorf("validate profile from env var: %v", err)
	}
	return profile, nil
}

// validateProfile validates the profile name.
// It can only use letters, numbers, or "-".
// If the profile is invalid, it returns an error.
func validateProfile(profile string) error {
	if profile == "" {
		// Not actually needed as the regext would catch this, but will provide a better error message
		return fmt.Errorf("profile name cannot be empty")
	}
	match, err := regexp.MatchString("^[a-zA-Z0-9-]+$", profile)
	if err != nil {
		return fmt.Errorf("match string regex: %v", err)
	}
	if !match {
		return fmt.Errorf("profile name can only contain letters, numbers, and \"-\"")
	}
	return nil
}
