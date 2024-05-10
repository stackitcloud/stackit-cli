package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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
		contents, exists, err := readFileIfExists(profileFilePath)
		if err != nil {
			return "", fmt.Errorf("read profile from file: %w", err)
		}
		if !exists {
			return "", nil
		}
		profile = contents
	}

	err := ValidateProfile(profile)
	if err != nil {
		return "", fmt.Errorf("validate profile: %w", err)
	}
	return profile, nil
}

// SetProfile sets the profile to be used by the CLI.
func SetProfile(profile string) error {
	err := ValidateProfile(profile)
	if err != nil {
		return fmt.Errorf("validate profile: %w", err)
	}

	profileFilePath := filepath.Join(configFolderPath, fmt.Sprintf("%s.%s", profileFileName, profileFileExtension))
	err = os.WriteFile(profileFilePath, []byte(profile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("write profile to file: %w", err)
	}
	return nil
}

// UnsetProfile removes the profile file.
// If the profile file does not exist, it does nothing.
func UnsetProfile() error {
	err := os.Remove(profileFilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove profile file: %w", err)
	}
	return nil
}

// ValidateProfile validates the profile name.
// It can only use letters, numbers, or "-" and cannot be empty.
// If the profile is invalid, it returns an error.
func ValidateProfile(profile string) error {
	match, err := regexp.MatchString("^[a-zA-Z0-9-]+$", profile)
	if err != nil {
		return fmt.Errorf("match string regex: %w", err)
	}
	if !match {
		return &errors.InvalidProfileNameError{
			Profile: profile,
		}
	}
	return nil
}
