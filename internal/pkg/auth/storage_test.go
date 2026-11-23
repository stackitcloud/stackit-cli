package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

func TestSetGetAuthField(t *testing.T) {
	var (
		testField1 authFieldKey = "test-field-1"
		testField2 authFieldKey = "test-field-2"
	)

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

func TestSetGetAuthFieldWithProfile(t *testing.T) {
	var (
		testField1 authFieldKey = "test-field-1"
		testField2 authFieldKey = "test-field-2"
	)

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
		activeProfile    string
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
			activeProfile: "test-profile",
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
			activeProfile: "test-profile",
			expectedValues: map[authFieldKey]string{
				testField1: testValue1,
				testField2: testValue2,
			},
		},
		{
			description:  "overlapping assignments",
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
			activeProfile: "test-profile",
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
			activeProfile: "test-profile",
			expectedValues: map[authFieldKey]string{
				testField1: testValue3,
				testField2: testValue2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Append random string to profile name to avoid conflicts
			tt.activeProfile = makeProfileNameUnique(tt.activeProfile)

			// Make sure profile name is valid
			err := config.ValidateProfile(tt.activeProfile)
			if err != nil {
				t.Fatalf("Profile name \"%s\" is invalid: %v", tt.activeProfile, err)
			}

			if !tt.keyringFails {
				keyring.MockInit()
			} else {
				keyring.MockInitWithError(fmt.Errorf("keyring unavailable for testing"))
			}

			for _, assignment := range tt.valueAssignments {
				err := setAuthFieldWithProfile(tt.activeProfile, assignment.key, assignment.value)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", assignment.key, assignment.value, err)
				}
				// Check that this value will be checked
				if _, ok := tt.expectedValues[assignment.key]; !ok {
					t.Fatalf("Value \"%s\" set but not checked. Please add it to 'expectedValues'", assignment.key)
				}
			}

			for key, valueExpected := range tt.expectedValues {
				value, err := getAuthFieldWithProfile(tt.activeProfile, key)
				if err != nil {
					t.Errorf("Failed to get value of \"%s\": %v", key, err)
					continue
				} else if value != valueExpected {
					t.Errorf("Value of field \"%s\" is wrong: expected \"%s\", got \"%s\"", key, valueExpected, value)
				}

				if !tt.keyringFails {
					err = deleteAuthFieldInKeyring(tt.activeProfile, key)
					if err != nil {
						t.Errorf("Post-test cleanup failed: remove field \"%s\" from keyring: %v. Please remove it manually", key, err)
					}
				} else {
					err = deleteAuthFieldInEncodedTextFile(tt.activeProfile, key)
					if err != nil {
						t.Errorf("Post-test cleanup failed: remove field \"%s\" from text file: %v. Please remove it manually", key, err)
					}
				}
			}

			err = deleteProfileFiles(tt.activeProfile)
			if err != nil {
				t.Errorf("Post-test cleanup failed: remove profile \"%s\": %v. Please remove it manually", tt.activeProfile, err)
			}
		})
	}
}

func TestSetGetAuthFieldKeyring(t *testing.T) {
	var (
		testField1 authFieldKey = "test-field-1"
		testField2 authFieldKey = "test-field-2"
	)

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

			// Apppend random string to profile name to avoid conflicts
			tt.activeProfile = makeProfileNameUnique(tt.activeProfile)

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

func TestDeleteAuthField(t *testing.T) {
	tests := []struct {
		description  string
		keyringFails bool
		key          authFieldKey
		noKey        bool
	}{
		{
			description: "base",
			key:         "test-field-1",
		},
		{
			description: "key doesnt exist",
			key:         "doesnt-exist",
			noKey:       true,
		},
		{
			description:  "keyring fails",
			keyringFails: true,
			key:          "test-field-1",
		},
		{
			description:  "keyring fails, no key exists",
			keyringFails: true,
			noKey:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if !tt.keyringFails {
				keyring.MockInit()
			} else {
				keyring.MockInitWithError(fmt.Errorf("keyring unavailable for testing"))
			}

			// Append random string to auth field key and value to avoid conflicts
			testField1 := authFieldKey(fmt.Sprintf("test-field-1-%s", time.Now().Format(time.RFC3339)))
			testValue1 := fmt.Sprintf("value-1-%s", time.Now().Format(time.RFC3339))

			if !tt.noKey {
				err := SetAuthField(testField1, testValue1)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", testField1, testValue1, err)
				}
			}

			err := DeleteAuthField(tt.key)
			if err != nil {
				t.Fatalf("Failed to delete field \"%s\": %v", tt.key, err)
			}

			// Check if key still exists
			_, err = GetAuthField(tt.key)
			if err == nil {
				t.Fatalf("Key \"%s\" still exists after deletion", tt.key)
			}
		})
	}
}

func TestDeleteAuthFieldWithProfile(t *testing.T) {
	tests := []struct {
		description  string
		keyringFails bool
		profile      string
		key          authFieldKey
		noKey        bool
	}{
		{
			description: "base",
			profile:     "default",
			key:         "test-field-1",
		},
		{
			description: "key doesnt exist",
			profile:     "default",
			key:         "doesnt-exist",
			noKey:       true,
		},
		{
			description:  "keyring fails",
			profile:      "default",
			keyringFails: true,
			key:          "test-field-1",
		},
		{
			description:  "keyring fails, no key exists",
			profile:      "default",
			keyringFails: true,
			noKey:        true,
		},
		{
			description: "base, custom profile",
			profile:     "test-profile",
			key:         "test-field-1",
		},
		{
			description: "key doesnt exist, custom profile",
			profile:     "test-profile",
			key:         "doesnt-exist",
			noKey:       true,
		},
		{
			description:  "keyring fails, custom profile",
			profile:      "test-profile",
			keyringFails: true,
			key:          "test-field-1",
		},
		{
			description:  "keyring fails, no key exists, custom profile",
			profile:      "test-profile",
			keyringFails: true,
			noKey:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if !tt.keyringFails {
				keyring.MockInit()
			} else {
				keyring.MockInitWithError(fmt.Errorf("keyring unavailable for testing"))
			}

			// Append random string to auth field key and value to avoid conflicts
			testField1 := authFieldKey(fmt.Sprintf("test-field-1-%s", time.Now().Format(time.RFC3339)))
			testValue1 := fmt.Sprintf("value-1-%s", time.Now().Format(time.RFC3339))

			if !tt.noKey {
				err := SetAuthField(testField1, testValue1)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", testField1, testValue1, err)
				}
			}

			err := deleteAuthFieldWithProfile(tt.profile, tt.key)
			if err != nil {
				t.Fatalf("Failed to delete field \"%s\": %v", tt.key, err)
			}

			// Check if key still exists
			_, err = GetAuthField(tt.key)
			if err == nil {
				t.Fatalf("Key \"%s\" still exists after deletion", tt.key)
			}
		})
	}
}

func TestDeleteAuthFieldKeyring(t *testing.T) {
	tests := []struct {
		description   string
		activeProfile string
		noKey         bool
		isValid       bool
	}{
		{
			description:   "base, default profile",
			activeProfile: config.DefaultProfileName,
			isValid:       true,
		},
		{
			description:   "key doesnt exist, default profile",
			activeProfile: config.DefaultProfileName,
			noKey:         true,
			isValid:       false,
		},
		{
			description:   "base, custom profile",
			activeProfile: "test-profile",
			isValid:       true,
		},
		{
			description:   "key doesnt exist, custom profile",
			activeProfile: "test-profile",
			noKey:         true,
			isValid:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()

			// Append random string to auth field key and value to avoid conflicts
			testField1 := authFieldKey(fmt.Sprintf("test-field-1-%s", time.Now().Format(time.RFC3339)))
			testValue1 := fmt.Sprintf("value-1-keyring-%s", time.Now().Format(time.RFC3339))

			// Append random string to profile name to avoid conflicts
			tt.activeProfile = makeProfileNameUnique(tt.activeProfile)

			// Make sure profile name is valid
			err := config.ValidateProfile(tt.activeProfile)
			if err != nil {
				t.Fatalf("Profile name \"%s\" is invalid: %v", tt.activeProfile, err)
			}

			if !tt.noKey {
				err := setAuthFieldInKeyring(tt.activeProfile, testField1, testValue1)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", testField1, testValue1, err)
				}
			}

			err = deleteAuthFieldInKeyring(tt.activeProfile, testField1)
			if err != nil {
				if tt.isValid {
					t.Fatalf("Failed to delete field \"%s\" from keyring: %v", testField1, err)
				}
				return
			}

			if !tt.isValid {
				t.Fatalf("Expected error when deleting field \"%s\" from keyring, got none", testField1)
			}

			// Check if key still exists
			_, err = getAuthFieldFromKeyring(tt.activeProfile, testField1)
			if err == nil {
				t.Fatalf("Key \"%s\" still exists in keyring after deletion", testField1)
			}
		})
	}
}

func TestDeleteProfileFromKeyring(t *testing.T) {
	tests := []struct {
		description   string
		keyringFails  bool
		keys          []authFieldKey
		activeProfile string
		isValid       bool
	}{

		{
			description:   "base",
			keys:          authFieldKeys,
			activeProfile: "test-profile",
			isValid:       true,
		},
		{
			description: "missing keys",
			keys: []authFieldKey{
				ACCESS_TOKEN,
				SERVICE_ACCOUNT_EMAIL,
			},
			activeProfile: "test-profile",
			isValid:       true,
		},
		{
			description:   "invalid profile",
			activeProfile: "INVALID",
			isValid:       false,
		},
		{
			description:  "keyring fails",
			keyringFails: true,
			isValid:      false,
		},
		{
			description:   "default profile",
			activeProfile: config.DefaultProfileName,
			isValid:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if !tt.keyringFails {
				keyring.MockInit()
			} else {
				keyring.MockInitWithError(fmt.Errorf("keyring unavailable for testing"))
			}

			// Append random string to auth field key and value to avoid conflicts
			testValue1 := fmt.Sprintf("value-1-keyring-%s", time.Now().Format(time.RFC3339))

			// Append random string to profile name to avoid conflicts
			tt.activeProfile = makeProfileNameUnique(tt.activeProfile)

			for _, key := range tt.keys {
				err := setAuthFieldInKeyring(tt.activeProfile, key, testValue1)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", key, testValue1, err)
				}
			}

			err := DeleteProfileAuth(tt.activeProfile)
			if err != nil {
				if tt.isValid {
					t.Fatalf("Failed to delete profile \"%s\" from keyring: %v", tt.activeProfile, err)
				}
				return
			}

			if !tt.isValid {
				t.Fatalf("Expected error when deleting profile \"%s\" from keyring, got none", tt.activeProfile)
			}

			for _, key := range tt.keys {
				// Check if key still exists
				_, err = getAuthFieldFromKeyring(tt.activeProfile, key)
				if err == nil {
					t.Fatalf("Key \"%s\" still exists in keyring after profile deletion", key)
				}
			}
		})
	}
}

func TestSetGetAuthFieldEncodedTextFile(t *testing.T) {
	var (
		testField1 authFieldKey = "test-field-1"
		testField2 authFieldKey = "test-field-2"
	)

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
			// Append random string to profile name to avoid conflicts
			tt.activeProfile = makeProfileNameUnique(tt.activeProfile)

			// Make sure profile name is valid
			err := config.ValidateProfile(tt.activeProfile)
			if err != nil {
				t.Fatalf("Profile name \"%s\" is invalid: %v", tt.activeProfile, err)
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

			err = deleteProfileFiles(tt.activeProfile)
			if err != nil {
				t.Errorf("Post-test cleanup failed: remove profile \"%s\": %v. Please remove it manually", tt.activeProfile, err)
			}
		})
	}
}

func TestGetProfileEmail(t *testing.T) {
	tests := []struct {
		description     string
		activeProfile   string
		userEmail       string
		authFlow        AuthFlow
		serviceAccEmail string
		expectedEmail   string
	}{
		{
			description:   "default profile, user token",
			activeProfile: config.DefaultProfileName,
			userEmail:     "test@test.com",
			authFlow:      AUTH_FLOW_USER_TOKEN,
			expectedEmail: "test@test.com",
		},
		{
			description:     "default profile, service acc token",
			activeProfile:   config.DefaultProfileName,
			serviceAccEmail: "test@test.com",
			authFlow:        AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			expectedEmail:   "test@test.com",
		},
		{
			description:     "default profile, service acc key",
			activeProfile:   config.DefaultProfileName,
			serviceAccEmail: "test@test.com",
			authFlow:        AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			expectedEmail:   "test@test.com",
		},
		{
			description:   "custom profile, user token",
			activeProfile: "test-profile",
			userEmail:     "test@test.com",
			authFlow:      AUTH_FLOW_USER_TOKEN,
			expectedEmail: "test@test.com",
		},
		{
			description:     "custom profile, service acc token",
			activeProfile:   "test-profile",
			serviceAccEmail: "test@test.com",
			authFlow:        AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			expectedEmail:   "test@test.com",
		},
		{
			description:     "custom profile, service acc key",
			activeProfile:   "test-profile",
			serviceAccEmail: "test@test.com",
			authFlow:        AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			expectedEmail:   "test@test.com",
		},
		{
			description:   "no email, user token",
			activeProfile: "test-profile",
			authFlow:      AUTH_FLOW_USER_TOKEN,
			expectedEmail: "",
		},
		{
			description:   "no email, service acc token",
			activeProfile: "test-profile",
			authFlow:      AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			expectedEmail: "",
		},
		{
			description:   "no email, service acc key",
			activeProfile: "test-profile",
			authFlow:      AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			expectedEmail: "",
		},
		{
			description:   "user not authenticated",
			activeProfile: "test-profile",
			expectedEmail: "",
		},
		{
			description:     "both emails, user not authenticated",
			activeProfile:   "test-profile",
			userEmail:       "test@test.com",
			serviceAccEmail: "test2@test.com",
			expectedEmail:   "",
		},
		{
			description:     "both emails, user token",
			activeProfile:   "test-profile",
			userEmail:       "test@test.com",
			serviceAccEmail: "test2@test.com",
			authFlow:        AUTH_FLOW_USER_TOKEN,
			expectedEmail:   "test@test.com",
		},
		{
			description:     "both emails, service account token",
			activeProfile:   "test-profile",
			userEmail:       "test@test.com",
			serviceAccEmail: "test2@test.com",
			authFlow:        AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			expectedEmail:   "test2@test.com",
		},
		{
			description:     "both emails, service account key",
			activeProfile:   "test-profile",
			userEmail:       "test@test.com",
			serviceAccEmail: "test2@test.com",
			authFlow:        AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			expectedEmail:   "test2@test.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()

			// Append random string to profile name to avoid conflicts
			tt.activeProfile = makeProfileNameUnique(tt.activeProfile)

			// Make sure profile name is valid
			err := config.ValidateProfile(tt.activeProfile)
			if err != nil {
				t.Fatalf("Profile name \"%s\" is invalid: %v", tt.activeProfile, err)
			}

			err = setAuthFieldInKeyring(tt.activeProfile, USER_EMAIL, tt.userEmail)
			if err != nil {
				t.Errorf("Failed to set user email: %v", err)
			}

			err = setAuthFieldInKeyring(tt.activeProfile, SERVICE_ACCOUNT_EMAIL, tt.serviceAccEmail)
			if err != nil {
				t.Errorf("Failed to set service account email: %v", err)
			}

			err = setAuthFieldWithProfile(tt.activeProfile, authFlowType, string(tt.authFlow))
			if err != nil {
				t.Errorf("Failed to set auth flow: %v", err)
			}

			email := GetProfileEmail(tt.activeProfile)
			if email != tt.expectedEmail {
				t.Errorf("Expected email \"%s\", got \"%s\"", tt.expectedEmail, email)
			}

			err = deleteAuthFieldInKeyring(tt.activeProfile, USER_EMAIL)
			if err != nil {
				t.Fatalf("Failed to remove user email: %v", err)
			}

			err = deleteAuthFieldInKeyring(tt.activeProfile, SERVICE_ACCOUNT_EMAIL)
			if err != nil {
				t.Fatalf("Failed to remove service account email: %v", err)
			}

			err = deleteProfileFiles(tt.activeProfile)
			if err != nil {
				t.Fatalf("Failed to remove profile: %v", err)
			}
		})
	}
}

func deleteProfileFiles(activeProfile string) error {
	if activeProfile == config.DefaultProfileName {
		// Do not delete the default profile
		return nil
	}

	textFileDir := config.GetProfileFolderPath(activeProfile)
	// Remove the entire directory if the profile does not exist
	err := os.RemoveAll(textFileDir)
	if err != nil {
		return fmt.Errorf("remove directory: %w", err)
	}

	return nil
}

func makeProfileNameUnique(profile string) string {
	if profile == config.DefaultProfileName {
		return profile
	}
	return fmt.Sprintf("%s-%s", profile, time.Now().Format("20060102150405"))
}

func TestAuthorizeDeauthorizeUserProfileAuth(t *testing.T) {
	type args struct {
		sessionExpiresAtUnix string
		accessToken          string
		refreshToken         string
		email                string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				sessionExpiresAtUnix: "1234567890",
				accessToken:          "accessToken",
				refreshToken:         "refreshToken",
				email:                "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "no email",
			args: args{
				sessionExpiresAtUnix: "1234567890",
				accessToken:          "accessToken",
				refreshToken:         "refreshToken",
				email:                "",
			},
			wantErr: false,
		},
		{
			name: "no session expires",
			args: args{
				sessionExpiresAtUnix: "",
				accessToken:          "accessToken",
				refreshToken:         "refreshToken",
				email:                "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "no access token",
			args: args{
				sessionExpiresAtUnix: "1234567890",
				accessToken:          "",
				refreshToken:         "refreshToken",
				email:                "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "no refresh token",
			args: args{
				sessionExpiresAtUnix: "1234567890",
				accessToken:          "accessToken",
				refreshToken:         "",
				email:                "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "all empty args",
			args: args{
				sessionExpiresAtUnix: "",
				accessToken:          "",
				refreshToken:         "",
				email:                "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyring.MockInit()

			if err := LoginUser(tt.args.email, tt.args.accessToken, tt.args.refreshToken, tt.args.sessionExpiresAtUnix); (err != nil) != tt.wantErr {
				t.Errorf("AuthorizeUserProfileAuth() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Test values
			testLoginAuthFields := []string{
				tt.args.sessionExpiresAtUnix,
				tt.args.accessToken,
				tt.args.refreshToken,
				tt.args.email,
			}

			// Check if the fields are set
			for i := range loginAuthFieldKeys {
				gotKey, err := GetAuthField(loginAuthFieldKeys[i])
				if err != nil {
					t.Errorf("Field \"%s\" not set after authorization", loginAuthFieldKeys[i])
				}

				expectedKey := testLoginAuthFields[i]
				if gotKey != expectedKey {
					t.Errorf("Field \"%s\" is wrong: expected \"%s\", got \"%s\"", loginAuthFieldKeys[i], expectedKey, gotKey)
				}
			}

			if err := LogoutUser(); err != nil {
				t.Errorf("DeauthorizeUserProfileAuth() error = %v", err)
			}

			// Check if the fields are deleted
			for _, key := range loginAuthFieldKeys {
				_, err := GetAuthField(key)
				if err == nil {
					t.Errorf("Field \"%s\" still exists after deauthorization", key)
				}
			}
		})
	}
}

func TestStaleKeychainCleanupOnFallback(t *testing.T) {
	keyring.MockInit()

	profile, err := config.GetProfile()
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}

	oldKeyringVal := "stale_rsa_2048_key"
	newFileVal := "new_rsa_4096_key"

	// 1. Initially set the old key in keyring
	err = setAuthFieldInKeyring(profile, PRIVATE_KEY, oldKeyringVal)
	if err != nil {
		t.Fatalf("set in keyring: %v", err)
	}

	// Verify keyring has it
	got, err := getAuthFieldFromKeyring(profile, PRIVATE_KEY)
	if err != nil || got != oldKeyringVal {
		t.Fatalf("expected %s in keyring, got %s (err: %v)", oldKeyringVal, got, err)
	}

	// 2. Simulate what setAuthFieldWithProfile does when keyring write fails:
	// It writes to encoded text file and deletes from keyring
	err = setAuthFieldInEncodedTextFile(profile, PRIVATE_KEY, newFileVal)
	if err != nil {
		t.Fatalf("set in text file: %v", err)
	}

	err = deleteAuthFieldInKeyring(profile, PRIVATE_KEY)
	if err != nil {
		t.Fatalf("delete from keyring: %v", err)
	}

	// 3. Now GetAuthField should return newFileVal, NOT the stale oldKeyringVal!
	val, err := GetAuthField(PRIVATE_KEY)
	if err != nil {
		t.Fatalf("GetAuthField failed: %v", err)
	}

	if val != newFileVal {
		t.Errorf("GetAuthField returned stale value %q, expected %q", val, newFileVal)
	}

	// Cleanup
	err = DeleteAuthField(PRIVATE_KEY)
	if err != nil {
		t.Fatalf("DeleteAuthField failed: %v", err)
	}
}

func TestFileNotCreatedOnReadOrDelete(t *testing.T) {
	keyring.MockInit()

	profile, err := config.GetProfile()
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}

	textFilePath := filepath.Join(config.GetProfileFolderPath(profile), textFileName)

	// Ensure file doesn't exist
	if _, err := os.Stat(textFilePath); !os.IsNotExist(err) {
		t.Fatalf("expected file to not exist before test")
	}

	nonExistingField := authFieldKey("non_existent_field_12345")

	// Call GetAuthField for a non-existent field
	_, err = GetAuthField(nonExistingField)
	if err == nil {
		t.Errorf("expected error from GetAuthField for non-existent field")
	}

	// Verify file was NOT created
	if _, err := os.Stat(textFilePath); !os.IsNotExist(err) {
		t.Errorf("cli-auth-storage.txt was created on GetAuthField for non-existent field")
	}

	// Call DeleteAuthField for a non-existent field
	_ = DeleteAuthField(nonExistingField)

	// Verify file was still NOT created
	if _, err := os.Stat(textFilePath); !os.IsNotExist(err) {
		t.Errorf("cli-auth-storage.txt was created on DeleteAuthField for non-existent field")
	}
}

func TestFileRemovedWhenEmpty(t *testing.T) {
	profile, err := config.GetProfile()
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}

	textFilePath := filepath.Join(config.GetProfileFolderPath(profile), textFileName)

	field1 := authFieldKey("test-cleanup-field-1")
	field2 := authFieldKey("test-cleanup-field-2")

	// Write two fields to text file
	if err := setAuthFieldInEncodedTextFile(profile, field1, "val1"); err != nil {
		t.Fatalf("failed to set field1: %v", err)
	}

	if err := setAuthFieldInEncodedTextFile(profile, field2, "val2"); err != nil {
		t.Fatalf("failed to set field2: %v", err)
	}

	if _, err := os.Stat(textFilePath); os.IsNotExist(err) {
		t.Fatalf("expected file to exist after writing fields")
	}

	// Delete field1
	if err := deleteAuthFieldInEncodedTextFile(profile, field1); err != nil {
		t.Fatalf("failed to delete field1: %v", err)
	}

	if _, err := os.Stat(textFilePath); os.IsNotExist(err) {
		t.Fatalf("expected file to still exist with field2 remaining")
	}

	// Delete field2
	if err := deleteAuthFieldInEncodedTextFile(profile, field2); err != nil {
		t.Fatalf("failed to delete field2: %v", err)
	}

	// Now file should be removed
	if _, err := os.Stat(textFilePath); !os.IsNotExist(err) {
		t.Errorf("expected cli-auth-storage.txt to be removed when empty, but it still exists")
	}
}

func TestLogoutCleansUpAllCredentialsAndFile(t *testing.T) {
	keyring.MockInit()

	profile, err := config.GetProfile()
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}

	textFilePath := filepath.Join(config.GetProfileFolderPath(profile), textFileName)

	// Set user credentials
	if err := SetAuthField(USER_EMAIL, "[EMAIL_REDACTED]"); err != nil {
		t.Fatalf("set user email: %v", err)
	}

	if err := SetAuthField(ACCESS_TOKEN, "user-access-token"); err != nil {
		t.Fatalf("set access token: %v", err)
	}

	// Set service account credentials
	if err := SetAuthField(SERVICE_ACCOUNT_EMAIL, "[EMAIL_REDACTED]"); err != nil {
		t.Fatalf("set sa email: %v", err)
	}

	if err := SetAuthField(SERVICE_ACCOUNT_KEY, `{"id":"test"}`); err != nil {
		t.Fatalf("set sa key: %v", err)
	}

	// Also simulate fallback text file having a field
	if err := setAuthFieldInEncodedTextFile(profile, PRIVATE_KEY, "some-private-key"); err != nil {
		t.Fatalf("set private key in text file: %v", err)
	}

	// Ensure text file exists
	if _, err := os.Stat(textFilePath); os.IsNotExist(err) {
		t.Fatalf("expected text file to exist")
	}

	// Perform logout
	if err := LogoutUser(); err != nil {
		t.Fatalf("LogoutUser() failed: %v", err)
	}

	// Verify user fields are deleted
	if _, err := GetAuthField(USER_EMAIL); err == nil {
		t.Errorf("USER_EMAIL still exists after logout")
	}

	if _, err := GetAuthField(ACCESS_TOKEN); err == nil {
		t.Errorf("ACCESS_TOKEN still exists after logout")
	}

	// Verify service account fields are deleted
	if _, err := GetAuthField(SERVICE_ACCOUNT_EMAIL); err == nil {
		t.Errorf("SERVICE_ACCOUNT_EMAIL still exists after logout")
	}

	if _, err := GetAuthField(SERVICE_ACCOUNT_KEY); err == nil {
		t.Errorf("SERVICE_ACCOUNT_KEY still exists after logout")
	}

	if _, err := GetAuthField(PRIVATE_KEY); err == nil {
		t.Errorf("PRIVATE_KEY still exists after logout")
	}

	// Verify text file was removed
	if _, err := os.Stat(textFilePath); !os.IsNotExist(err) {
		t.Errorf("cli-auth-storage.txt was not removed after logout")
	}
}

func TestLoginUserCleansUpServiceAccountCredentials(t *testing.T) {
	keyring.MockInit()

	// Set SA credentials
	if err := SetAuthField(SERVICE_ACCOUNT_EMAIL, "[EMAIL_REDACTED]"); err != nil {
		t.Fatalf("set sa email: %v", err)
	}

	if err := SetAuthField(SERVICE_ACCOUNT_KEY, `{"id":"test"}`); err != nil {
		t.Fatalf("set sa key: %v", err)
	}

	if err := SetAuthField(PRIVATE_KEY, "private-key-data"); err != nil {
		t.Fatalf("set private key: %v", err)
	}

	// Now call LoginUser
	if err := LoginUser("[EMAIL_REDACTED]", "user-token", "refresh-token", "1234567890"); err != nil {
		t.Fatalf("LoginUser failed: %v", err)
	}

	// Verify SA fields were cleaned up
	if val, err := GetAuthField(SERVICE_ACCOUNT_EMAIL); err == nil && val != "" {
		t.Errorf("SERVICE_ACCOUNT_EMAIL still exists after LoginUser: %s", val)
	}

	if val, err := GetAuthField(SERVICE_ACCOUNT_KEY); err == nil && val != "" {
		t.Errorf("SERVICE_ACCOUNT_KEY still exists after LoginUser: %s", val)
	}

	if val, err := GetAuthField(PRIVATE_KEY); err == nil && val != "" {
		t.Errorf("PRIVATE_KEY still exists after LoginUser: %s", val)
	}

	// Verify user fields are present
	if val, err := GetAuthField(USER_EMAIL); err != nil || val != "[EMAIL_REDACTED]" {
		t.Errorf("expected USER_EMAIL '[EMAIL_REDACTED]', got %q (err: %v)", val, err)
	}

	// Cleanup
	err := LogoutUser()
	if err != nil {
		t.Fatalf("LogoutUser() failed: %v", err)
	}
}

func TestSetAuthFieldWithProfile(t *testing.T) {
	tests := []struct {
		description       string
		keyringMockErr    error
		setupCorruptFile  bool
		setupStaleFileVal string
		key               authFieldKey
		value             string
		wantErr           bool
		errContains       string
		verifyInKeyring   bool
		verifyInFile      bool
	}{
		{
			description:       "keyring write succeeds and cleans up stale text file entry",
			setupStaleFileVal: "old-file-val",
			key:               "test-sync-key",
			value:             "new-keyring-val",
			wantErr:           false,
			verifyInKeyring:   true,
		},
		{
			description:      "keyring write succeeds but fallback file cleanup fails",
			setupCorruptFile: true,
			key:              "test-corrupt-file-key",
			value:            "value",
			wantErr:          true,
			errContains:      "saved to keyring, but failed to clean up legacy fallback file",
		},
		{
			description:    "keyring write fails and keyring cleanup fails with unexpected error",
			keyringMockErr: fmt.Errorf("access denied to secure enclave"),
			key:            "test-keyring-delete-fail",
			value:          "value",
			wantErr:        true,
			errContains:    "cleanup stale key from keyring",
		},
		{
			description:    "keyring write fails and keyring cleanup is ignored when keyring is unavailable",
			keyringMockErr: fmt.Errorf("dbus daemon unavailable"),
			key:            "test-keyring-unavailable",
			value:          "value",
			wantErr:        false,
			verifyInFile:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			profile, err := config.GetProfile()
			if err != nil {
				t.Fatalf("get profile: %v", err)
			}

			if tt.keyringMockErr != nil {
				keyring.MockInitWithError(tt.keyringMockErr)
			} else {
				keyring.MockInit()
			}

			textFilePath := filepath.Join(config.GetProfileFolderPath(profile), textFileName)
			if tt.setupCorruptFile {
				if err := os.WriteFile(textFilePath, []byte("invalid-base64!!!"), 0o600); err != nil {
					t.Fatalf("failed to write corrupt file: %v", err)
				}
				defer func() {
					_ = os.Remove(textFilePath)
				}()
			} else if tt.setupStaleFileVal != "" {
				if err := setAuthFieldInEncodedTextFile(profile, tt.key, tt.setupStaleFileVal); err != nil {
					t.Fatalf("failed to set stale file value: %v", err)
				}
			}

			defer func() {
				_ = deleteAuthFieldInEncodedTextFile(profile, tt.key)
				_ = deleteAuthFieldInKeyring(profile, tt.key)
			}()

			err = setAuthFieldWithProfile(profile, tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("setAuthFieldWithProfile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("expected error containing %q, got: %v", tt.errContains, err)
			}

			if tt.verifyInKeyring {
				got, err := getAuthFieldFromKeyring(profile, tt.key)
				if err != nil || got != tt.value {
					t.Errorf("expected %q in keyring, got %q (err: %v)", tt.value, got, err)
				}

				if _, err := getAuthFieldFromEncodedTextFile(profile, tt.key); err == nil {
					t.Errorf("expected key to be removed from text file")
				}
			}

			if tt.verifyInFile {
				got, err := getAuthFieldFromEncodedTextFile(profile, tt.key)
				if err != nil || got != tt.value {
					t.Errorf("expected %q in text file, got %q (err: %v)", tt.value, got, err)
				}
			}
		})
	}
}
