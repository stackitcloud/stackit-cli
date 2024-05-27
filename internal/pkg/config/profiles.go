package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const ProfileEnvVar = "STACKIT_CLI_PROFILE"

// GetProfile returns the current profile to be used by the CLI.
//
// The profile is determined by the value of the STACKIT_CLI_PROFILE environment variable, or, if not set,
// by the contents of the profile file in the CLI config folder.
//
// If the environment variable is not set and the profile file does not exist, it returns an empty string.
//
// If the profile is not valid, it returns an error.
func GetProfile() (string, error) {
	profile, profileSet := os.LookupEnv(ProfileEnvVar)
	if !profileSet {
		contents, exists, err := fileutils.ReadFileIfExists(profileFilePath)
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

// GetProfileFromEnv returns the profile from the environment variable.
// If the environment variable is not set, it returns an empty string.
// If the profile is not valid, it returns an error.
func GetProfileFromEnv() (string, bool) {
	return os.LookupEnv(ProfileEnvVar)
}

// CreateProfile creates a new profile.
// It creates a new folder for the profile and copies the config file from the current profile to the new profile.
// If the fromDefault flag is set, it creates an empty profile.
// If the setProfile flag is set, it sets the new profile as the active profile.
// If the profile already exists, it returns an error.
func CreateProfile(p *print.Printer, profile string, setProfile, emptyProfile bool) error {
	err := ValidateProfile(profile)
	if err != nil {
		return fmt.Errorf("validate profile: %w", err)
	}

	configFolderPath = filepath.Join(defaultConfigFolderPath, profileRootFolder, profile)

	// Error if the profile already exists
	_, err = os.Stat(configFolderPath)
	if err == nil {
		return fmt.Errorf("profile %q already exists", profile)
	}

	err = os.MkdirAll(configFolderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create config folder: %w", err)
	}
	p.Debug(print.DebugLevel, "created folder for the new profile: %s", configFolderPath)

	currentProfile, err := GetProfile()
	if err != nil {
		// Cleanup created directory
		cleanupErr := os.RemoveAll(configFolderPath)
		if cleanupErr != nil {
			return fmt.Errorf("get active profile: %w, cleanup directories: %w", err, cleanupErr)
		}
		return fmt.Errorf("get active profile: %w", err)
	}

	p.Debug(print.DebugLevel, "current active profile: %q", currentProfile)

	if !emptyProfile {
		p.Debug(print.DebugLevel, "duplicating profile configuration from %q to new profile %q", currentProfile, profile)
		err = DuplicateProfileConfiguration(p, currentProfile, profile)
		if err != nil {
			// Cleanup created directory
			cleanupErr := os.RemoveAll(configFolderPath)
			if cleanupErr != nil {
				return fmt.Errorf("get active profile: %w, cleanup directories: %w", err, cleanupErr)
			}
			return fmt.Errorf("duplicate profile configuration: %w", err)
		}
	}

	if setProfile {
		err = SetProfile(p, profile)
		if err != nil {
			return fmt.Errorf("set profile: %w", err)
		}
	}

	return nil
}

// DuplicateProfileConfiguration duplicates the current profile configuration to a new profile.
// It copies the config file from the current profile to the new profile.
// If the current profile does not exist, it returns an error.
// If the new profile already exists, it will be overwritten.
func DuplicateProfileConfiguration(p *print.Printer, currentProfile, newProfile string) error {
	var currentConfigFilePath string
	// If the current profile is empty, its the default profile
	if currentProfile == "" {
		currentConfigFilePath = filepath.Join(defaultConfigFolderPath, fmt.Sprintf("%s.%s", configFileName, configFileExtension))
	} else {
		currentConfigFilePath = filepath.Join(defaultConfigFolderPath, profileRootFolder, currentProfile, fmt.Sprintf("%s.%s", configFileName, configFileExtension))
	}

	newConfigFilePath := filepath.Join(configFolderPath, fmt.Sprintf("%s.%s", configFileName, configFileExtension))

	err := fileutils.CopyFile(currentConfigFilePath, newConfigFilePath)
	if err != nil {
		return fmt.Errorf("copy config file: %w", err)
	}

	p.Debug(print.DebugLevel, "created new configuration for profile %q based on %q in: %s", newProfile, currentProfile, newConfigFilePath)

	return nil
}

// SetProfile sets the profile to be used by the CLI.
func SetProfile(p *print.Printer, profile string) error {
	err := ValidateProfile(profile)
	if err != nil {
		return fmt.Errorf("validate profile: %w", err)
	}

	err = os.WriteFile(profileFilePath, []byte(profile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("write profile to file: %w", err)
	}
	p.Debug(print.DebugLevel, "persisted new active profile in: %s", profileFilePath)

	configFolderPath = filepath.Join(defaultConfigFolderPath, profile)
	p.Debug(print.DebugLevel, "profile %q is now active", profile)

	return nil
}

// UnsetProfile removes the profile file.
// If the profile file does not exist, it does nothing.
func UnsetProfile(p *print.Printer) error {
	err := os.Remove(profileFilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove profile file: %w", err)
	}
	p.Debug(print.DebugLevel, "removed active profile file: %s", profileFilePath)
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
