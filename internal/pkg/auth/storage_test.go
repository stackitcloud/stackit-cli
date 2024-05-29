package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

func TestSetGetAuthField(t *testing.T) {
	var testField1 authFieldKey = "test-field-1"
	var testField2 authFieldKey = "test-field-2"

	testValue1 := fmt.Sprintf("value-1-%s", time.Now().Format(time.RFC3339))
	testValue2 := fmt.Sprintf("value-2-%s", time.Now().Format(time.RFC3339))
	testValue3 := fmt.Sprintf("value-3-%s", time.Now().Format(time.RFC3339))

	type valueAssignment struct {
		key   authFieldKey
		value string
	}

	tests := []struct {
		description      string
		keyringFails     bool
		valueAssignments []valueAssignment
		expectedValues   map[authFieldKey]string
	}{
		{
			description: "simple assignments",
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue1,
				testField2: testValue2,
			},
		},
		{
			description:  "simple assignments w/ keyring failing",
			keyringFails: true,
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue1,
				testField2: testValue2,
			},
		},
		{
			description: "overlapping assignments",
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
				{
					key:   testField1,
					value: testValue3,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue3,
				testField2: testValue2,
			},
		},
		{
			description:  "overlapping assignments w/ keyring failing",
			keyringFails: true,
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
				{
					key:   testField1,
					value: testValue3,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue3,
				testField2: testValue2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			activeProfile, err := config.GetProfile()
			if err != nil {
				t.Errorf("get profile: %v", err)
			}

			if !tt.keyringFails {
				keyring.MockInit()
			} else {
				keyring.MockInitWithError(fmt.Errorf("keyring unavailable for testing"))
			}

			for _, assignment := range tt.valueAssignments {
				err := SetAuthField(assignment.key, assignment.value)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", assignment.key, assignment.value, err)
				}
				// Check that this value will be checked
				if _, ok := tt.expectedValues[assignment.key]; !ok {
					t.Fatalf("Value \"%s\" set but not checked. Please add it to 'expectedValues'", assignment.key)
				}
			}

			for key, valueExpected := range tt.expectedValues {
				value, err := GetAuthField(key)
				if err != nil {
					t.Errorf("Failed to get value of \"%s\": %v", key, err)
					continue
				} else if value != valueExpected {
					t.Errorf("Value of field \"%s\" is wrong: expected \"%s\", got \"%s\"", key, valueExpected, value)
				}

				if !tt.keyringFails {
					err = deleteAuthFieldInKeyring(activeProfile, key)
					if err != nil {
						t.Errorf("Post-test cleanup failed: remove field \"%s\" from keyring: %v. Please remove it manually", key, err)
					}
				} else {
					err = deleteAuthFieldInEncodedTextFile(activeProfile, key)
					if err != nil {
						t.Errorf("Post-test cleanup failed: remove field \"%s\" from text file: %v. Please remove it manually", key, err)
					}
				}
			}
		})
	}
}

func TestSetGetAuthFieldKeyring(t *testing.T) {
	var testField1 authFieldKey = "test-field-1"
	var testField2 authFieldKey = "test-field-2"

	testValue1 := fmt.Sprintf("value-1-keyring-%s", time.Now().Format(time.RFC3339))
	testValue2 := fmt.Sprintf("value-2-keyring-%s", time.Now().Format(time.RFC3339))
	testValue3 := fmt.Sprintf("value-3-keyring-%s", time.Now().Format(time.RFC3339))

	type valueAssignment struct {
		key   authFieldKey
		value string
	}

	tests := []struct {
		description      string
		valueAssignments []valueAssignment
		expectedValues   map[authFieldKey]string
		activeProfile    string
	}{
		{
			description:   "simple assignments with default profile",
			activeProfile: config.DefaultProfileName,
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue1,
				testField2: testValue2,
			},
		},
		{
			description:   "overlapping assignments with default profile",
			activeProfile: config.DefaultProfileName,
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
				{
					key:   testField1,
					value: testValue3,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue3,
				testField2: testValue2,
			},
		},
		{
			description:   "simple assignments with test-profile",
			activeProfile: "test-profile",
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue1,
				testField2: testValue2,
			},
		},
		{
			description:   "overlapping assignments with test-profile",
			activeProfile: "test-profile",
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
				{
					key:   testField1,
					value: testValue3,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue3,
				testField2: testValue2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()

			// Make sure profile name is valid
			err := config.ValidateProfile(tt.activeProfile)
			if err != nil {
				t.Fatalf("Profile name \"%s\" is invalid: %v", tt.activeProfile, err)
			}

			for _, assignment := range tt.valueAssignments {
				err := setAuthFieldInKeyring(tt.activeProfile, assignment.key, assignment.value)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", assignment.key, assignment.value, err)
				}
				// Check that this value will be checked
				if _, ok := tt.expectedValues[assignment.key]; !ok {
					t.Fatalf("Value \"%s\" set but not checked. Please add it to 'expectedValues'", assignment.key)
				}
			}

			for key, valueExpected := range tt.expectedValues {
				value, err := getAuthFieldFromKeyring(tt.activeProfile, key)
				if err != nil {
					t.Errorf("Failed to get value of \"%s\": %v", key, err)
					continue
				} else if value != valueExpected {
					t.Errorf("Value of field \"%s\" is wrong: expected \"%s\", got \"%s\"", key, valueExpected, value)
				}

				err = deleteAuthFieldInKeyring(tt.activeProfile, key)
				if err != nil {
					t.Errorf("Post-test cleanup failed: remove field \"%s\" from keyring: %v. Please remove it manually", key, err)
				}
			}
		})
	}
}

func TestSetGetAuthFieldEncodedTextFile(t *testing.T) {
	var testField1 authFieldKey = "test-field-1"
	var testField2 authFieldKey = "test-field-2"

	testValue1 := fmt.Sprintf("value-1-text-%s", time.Now().Format(time.RFC3339))
	testValue2 := fmt.Sprintf("value-2-text-%s", time.Now().Format(time.RFC3339))
	testValue3 := fmt.Sprintf("value-3-text-%s", time.Now().Format(time.RFC3339))

	type valueAssignment struct {
		key   authFieldKey
		value string
	}

	tests := []struct {
		description      string
		activeProfile    string
		valueAssignments []valueAssignment
		expectedValues   map[authFieldKey]string
	}{
		{
			description:   "simple assignments with default profile",
			activeProfile: config.DefaultProfileName,
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue1,
				testField2: testValue2,
			},
		},
		{
			description:   "overlapping assignments with default profile",
			activeProfile: config.DefaultProfileName,
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
				{
					key:   testField1,
					value: testValue3,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue3,
				testField2: testValue2,
			},
		},
		{
			description:   "simple assignments with test-profile",
			activeProfile: "test-profile",
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue1,
				testField2: testValue2,
			},
		},
		{
			description:   "overlapping assignments with test-profile",
			activeProfile: "test-profile",
			valueAssignments: []valueAssignment{
				{
					key:   testField1,
					value: testValue1,
				},
				{
					key:   testField2,
					value: testValue2,
				},
				{
					key:   testField1,
					value: testValue3,
				},
			},
			expectedValues: map[authFieldKey]string{
				testField1: testValue3,
				testField2: testValue2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Make sure profile name is valid
			err := config.ValidateProfile(tt.activeProfile)
			if err != nil {
				t.Fatalf("Profile name \"%s\" is invalid: %v", tt.activeProfile, err)
			}

			// Check if the profile existed before the test, and if it didn't, delete it after the test
			profileExists, err := config.ProfileExists(tt.activeProfile)
			if err != nil {
				t.Fatalf("Failed to check if profile exists: %v", err)
			}

			for _, assignment := range tt.valueAssignments {
				err := setAuthFieldInEncodedTextFile(tt.activeProfile, assignment.key, assignment.value)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", assignment.key, assignment.value, err)
				}
				// Check that this value will be checked
				if _, ok := tt.expectedValues[assignment.key]; !ok {
					t.Fatalf("Value \"%s\" set but not checked. Please add it to 'expectedValues'", assignment.key)
				}
			}

			for key, valueExpected := range tt.expectedValues {
				value, err := getAuthFieldFromEncodedTextFile(tt.activeProfile, key)
				if err != nil {
					t.Errorf("Failed to get value of \"%s\": %v", key, err)
					continue
				} else if value != valueExpected {
					t.Errorf("Value of field \"%s\" is wrong: expected \"%s\", got \"%s\"", key, valueExpected, value)
				}

				err = deleteAuthFieldInEncodedTextFile(tt.activeProfile, key)
				if err != nil {
					t.Errorf("Post-test cleanup failed: remove field \"%s\" from text file: %v. Please remove it manually", key, err)
				}
			}

			err = deleteAuthFieldProfile(tt.activeProfile, profileExists)
			if err != nil {
				t.Errorf("Post-test cleanup failed: remove profile \"%s\": %v. Please remove it manually", tt.activeProfile, err)
			}
		})
	}
}

func TestGetProfileEmail(t *testing.T) {
	tests := []struct {
		description     string
		profile         string
		userEmail       string
		serviceAccEmail string
	}{
		{
			description: "default profile, user email",
			profile:     config.DefaultProfileName,
			userEmail:   "test@test.com",
		},
		{
			description:     "default profile, service acc email",
			profile:         config.DefaultProfileName,
			serviceAccEmail: "test@test.com",
		},
		{
			description: "custom profile, user email",
			profile:     "test-profile",
			userEmail:   "test@test.com",
		},
		{
			description:     "custom profile, service acc email",
			profile:         "test-profile",
			serviceAccEmail: "test@test.com",
		},
		{
			description: "none of the emails",
			profile:     "test-profile",
		},
		{
			description:     "both emails",
			profile:         "test-profile",
			userEmail:       "test@test.com",
			serviceAccEmail: "test2@test.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			profileExists, err := config.ProfileExists(tt.profile)
			if err != nil {
				t.Fatalf("Failed to check if profile exists: %v", err)
			}
			oldUserEmail, _ := getAuthFieldWithProfile(tt.profile, USER_EMAIL)
			oldServiceAccEmail, _ := getAuthFieldWithProfile(tt.profile, SERVICE_ACCOUNT_EMAIL)

			err = setAuthFieldInKeyring(tt.profile, USER_EMAIL, tt.userEmail)
			if err != nil {
				t.Errorf("Failed to set user email: %v", err)
			}

			err = setAuthFieldInKeyring(tt.profile, SERVICE_ACCOUNT_EMAIL, tt.serviceAccEmail)
			if err != nil {
				t.Errorf("Failed to set service account email: %v", err)
			}

			email := GetProfileEmail(tt.profile)
			if tt.userEmail != "" {
				if email != tt.userEmail {
					t.Errorf("User email is wrong: expected \"%s\", got \"%s\"", tt.userEmail, email)
				}
			} else if tt.serviceAccEmail != "" {
				if email != tt.serviceAccEmail {
					t.Errorf("Service account email is wrong: expected \"%s\", got \"%s\"", tt.serviceAccEmail, email)
				}
			} else {
				if email != "" {
					t.Errorf("Email is wrong: expected \"\", got \"%s\"", email)
				}
			}

			err = deleteAuthFieldInKeyring(tt.profile, USER_EMAIL)
			if err != nil {
				t.Fatalf("Failed to remove user email: %v", err)
			}

			err = deleteAuthFieldInKeyring(tt.profile, SERVICE_ACCOUNT_EMAIL)
			if err != nil {
				t.Fatalf("Failed to remove service account email: %v", err)
			}

			if oldUserEmail != "" {
				err := setAuthFieldInKeyring(tt.profile, USER_EMAIL, oldUserEmail)
				if err != nil {
					t.Fatalf("Failed to set back user email: %v", err)
				}
			}

			if oldServiceAccEmail != "" {
				err := setAuthFieldInKeyring(tt.profile, SERVICE_ACCOUNT_EMAIL, oldServiceAccEmail)
				if err != nil {
					t.Fatalf("Failed to set back service account email: %v", err)
				}
			}

			err = deleteAuthFieldProfile(tt.profile, profileExists)
			if err != nil {
				t.Fatalf("Failed to remove profile: %v", err)
			}
		})
	}
}

func deleteAuthFieldInKeyring(activeProfile string, key authFieldKey) error {
	if activeProfile != "" {
		activeProfileKeyring := filepath.Join(keyringService, activeProfile)
		return keyring.Delete(activeProfileKeyring, string(key))
	}

	return keyring.Delete(keyringService, string(key))
}

func deleteAuthFieldInEncodedTextFile(activeProfile string, key authFieldKey) error {
	err := createEncodedTextFile(activeProfile)
	if err != nil {
		return err
	}

	textFileDir := config.GetProfileFolderPath(activeProfile)
	textFilePath := filepath.Join(textFileDir, textFileName)

	contentEncoded, err := os.ReadFile(textFilePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	contentBytes, err := base64.StdEncoding.DecodeString(string(contentEncoded))
	if err != nil {
		return fmt.Errorf("decode file: %w", err)
	}
	content := map[authFieldKey]string{}
	err = json.Unmarshal(contentBytes, &content)
	if err != nil {
		return fmt.Errorf("unmarshal file: %w", err)
	}

	delete(content, key)

	contentBytes, err = json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshal file: %w", err)
	}
	contentEncoded = []byte(base64.StdEncoding.EncodeToString(contentBytes))
	err = os.WriteFile(textFilePath, contentEncoded, 0o600)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func deleteAuthFieldProfile(activeProfile string, profileExisted bool) error {
	textFileDir := config.GetProfileFolderPath(activeProfile)
	if !profileExisted {
		// Remove the entire directory if the profile does not exist
		err := os.RemoveAll(textFileDir)
		if err != nil {
			return fmt.Errorf("remove directory: %w", err)
		}
	}
	return nil
}
