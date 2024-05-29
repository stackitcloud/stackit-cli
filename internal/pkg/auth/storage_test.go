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
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
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
			description:   "simple assignments with testProfile",
			activeProfile: "testProfile",
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
			description:   "overlapping assignments with testProfile",
			activeProfile: "testProfile",
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
			description:   "simple assignments with testProfile",
			activeProfile: "testProfile",
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
			description:   "overlapping assignments with testProfile",
			activeProfile: "testProfile",
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

			// Create profile if it does not exist
			// Will be deleted at the end of the test
			profileExists, err := config.ProfileExists(tt.activeProfile)
			if err != nil {
				t.Fatalf("Failed to check if profile exists: %v", err)
			}
			if !profileExists {
				p := print.NewPrinter()
				err := config.CreateProfile(p, tt.activeProfile, true, true)
				if err != nil {
					t.Fatalf("Failed to create profile: %v", err)
				}
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
