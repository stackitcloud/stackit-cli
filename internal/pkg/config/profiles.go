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
// If the profile is not set (env var or profile file) or is set but does not exist, it falls back to the default profile.
//
// If the profile is not valid, it returns an error.
func GetProfile() (string, error) {
	profile, err := GetConfiguredProfile()
	if err != nil {
		return "", err
	}

	// Make sure the profile exists
	profileExists, err := ProfileExists(profile)
	if err != nil {
		return "", fmt.Errorf("check if profile exists: %w", err)
	}
	if !profileExists {
		return DefaultProfileName, nil
	}

	return profile, nil
}

// GetConfiguredProfile returns the profile that is configured by the user, which may not exist.
// The profile is determined by the value of the STACKIT_CLI_PROFILE environment variable, or, if not set,
// by the contents of the profile file in the CLI config folder.
//
// If the profile is not valid, it returns an error.
func GetConfiguredProfile() (string, error) {
	profile, profileSetInEnv := GetProfileFromEnv()
	if !profileSetInEnv {
		contents, exists, err := fileutils.ReadFileIfExists(profileFilePath)
		if err != nil {
			return "", fmt.Errorf("read profile from file: %w", err)
		}
		if !exists {
			return DefaultProfileName, nil
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
// If emptyProfile is true, it creates an empty profile. Otherwise, copies the config from the current profile to the new profile.
// If setProfile is true, it sets the new profile as the active profile.
// If the profile already exists, it returns an error.
func CreateProfile(p *print.Printer, profile string, setProfile, emptyProfile bool) error {
	err := ValidateProfile(profile)
	if err != nil {
		return fmt.Errorf("validate profile: %w", err)
	}

	// Cannot create a profile with the default name
	if profile == DefaultProfileName {
		return &errors.InvalidProfileNameError{
			Profile: profile,
		}
	}

	configFolderPath = GetProfileFolderPath(profile)

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

	if !emptyProfile {
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
	currentProfileFolder := GetProfileFolderPath(currentProfile)
	currentConfigFilePath := getConfigFilePath(currentProfileFolder)

	newConfigFilePath := getConfigFilePath(configFolderPath)

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

	profileExists, err := ProfileExists(profile)
	if err != nil {
		return fmt.Errorf("check if profile exists: %w", err)
	}

	if !profileExists {
		return &errors.SetInexistentProfile{Profile: profile}
	}

	err = os.WriteFile(profileFilePath, []byte(profile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("write profile to file: %w", err)
	}
	p.Debug(print.DebugLevel, "persisted new active profile in: %s", profileFilePath)

	configFolderPath = GetProfileFolderPath(profile)
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

func ProfileExists(profile string) (bool, error) {
	_, err := os.Stat(GetProfileFolderPath(profile))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("get profile folder: %w", err)
	}
	return true, nil
}

func GetProfileFolderPath(profile string) string {
	if defaultConfigFolderPath == "" {
		defaultConfigFolderPath = getInitialConfigDir()
	}

	if profile == DefaultProfileName {
		return defaultConfigFolderPath
	}
	return filepath.Join(defaultConfigFolderPath, profileRootFolder, profile)
}
