package auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/zalando/go-keyring"
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

func TestSetGetAuthFieldWithProfile(t *testing.T) {
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
			// Apppend random string to profile name to avoid conflicts
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

// TestStorageContextSeparation tests that CLI and Provider contexts use different keyring services
func TestStorageContextSeparation(t *testing.T) {
	var testField authFieldKey = "test-field-context"
	testValueCLI := fmt.Sprintf("cli-value-%s", time.Now().Format(time.RFC3339))
	testValueProvider := fmt.Sprintf("provider-value-%s", time.Now().Format(time.RFC3339))

	tests := []struct {
		description  string
		keyringFails bool
	}{
		{
			description: "with keyring",
		},
		{
			description:  "with file fallback",
			keyringFails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if !tt.keyringFails {
				keyring.MockInit()
			} else {
				keyring.MockInitWithError(fmt.Errorf("keyring unavailable for testing"))
			}

			// Set value in CLI context
			err := SetAuthFieldWithContext(StorageContextCLI, testField, testValueCLI)
			if err != nil {
				t.Fatalf("Failed to set CLI context field: %v", err)
			}

			// Set value in Provider context
			err = SetAuthFieldWithContext(StorageContextAPI, testField, testValueProvider)
			if err != nil {
				t.Fatalf("Failed to set Provider context field: %v", err)
			}

			// Verify CLI context value
			valueCLI, err := GetAuthFieldWithContext(StorageContextCLI, testField)
			if err != nil {
				t.Fatalf("Failed to get CLI context field: %v", err)
			}
			if valueCLI != testValueCLI {
				t.Errorf("CLI context value incorrect: expected %s, got %s", testValueCLI, valueCLI)
			}

			// Verify Provider context value
			valueProvider, err := GetAuthFieldWithContext(StorageContextAPI, testField)
			if err != nil {
				t.Fatalf("Failed to get Provider context field: %v", err)
			}
			if valueProvider != testValueProvider {
				t.Errorf("Provider context value incorrect: expected %s, got %s", testValueProvider, valueProvider)
			}

			// Cleanup
			activeProfile, _ := config.GetProfile()
			if !tt.keyringFails {
				_ = deleteAuthFieldInKeyringWithContext(StorageContextCLI, activeProfile, testField)
				_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, activeProfile, testField)
			} else {
				_ = deleteAuthFieldInEncodedTextFileWithContext(StorageContextCLI, activeProfile, testField)
				_ = deleteAuthFieldInEncodedTextFileWithContext(StorageContextAPI, activeProfile, testField)
			}
		})
	}
}

// TestStorageContextIsolation tests that changes in one context don't affect the other
func TestStorageContextIsolation(t *testing.T) {
	var testField authFieldKey = "test-field-isolation"
	testValueCLI := fmt.Sprintf("cli-value-%s", time.Now().Format(time.RFC3339))
	testValueProvider := fmt.Sprintf("provider-value-%s", time.Now().Format(time.RFC3339))
	updatedValueCLI := fmt.Sprintf("cli-updated-%s", time.Now().Format(time.RFC3339))

	keyring.MockInit()

	// Set values in both contexts
	err := SetAuthFieldWithContext(StorageContextCLI, testField, testValueCLI)
	if err != nil {
		t.Fatalf("Failed to set CLI context field: %v", err)
	}

	err = SetAuthFieldWithContext(StorageContextAPI, testField, testValueProvider)
	if err != nil {
		t.Fatalf("Failed to set Provider context field: %v", err)
	}

	// Update CLI context value
	err = SetAuthFieldWithContext(StorageContextCLI, testField, updatedValueCLI)
	if err != nil {
		t.Fatalf("Failed to update CLI context field: %v", err)
	}

	// Verify CLI context was updated
	valueCLI, err := GetAuthFieldWithContext(StorageContextCLI, testField)
	if err != nil {
		t.Fatalf("Failed to get CLI context field: %v", err)
	}
	if valueCLI != updatedValueCLI {
		t.Errorf("CLI context value not updated: expected %s, got %s", updatedValueCLI, valueCLI)
	}

	// Verify Provider context was NOT affected
	valueProvider, err := GetAuthFieldWithContext(StorageContextAPI, testField)
	if err != nil {
		t.Fatalf("Failed to get Provider context field: %v", err)
	}
	if valueProvider != testValueProvider {
		t.Errorf("Provider context value changed unexpectedly: expected %s, got %s", testValueProvider, valueProvider)
	}

	// Cleanup
	activeProfile, _ := config.GetProfile()
	_ = deleteAuthFieldInKeyringWithContext(StorageContextCLI, activeProfile, testField)
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, activeProfile, testField)
}

// TestStorageContextDeletion tests that deleting from one context doesn't affect the other
func TestStorageContextDeletion(t *testing.T) {
	var testField authFieldKey = "test-field-deletion"
	testValueCLI := fmt.Sprintf("cli-value-%s", time.Now().Format(time.RFC3339))
	testValueProvider := fmt.Sprintf("provider-value-%s", time.Now().Format(time.RFC3339))

	keyring.MockInit()

	// Set values in both contexts
	err := SetAuthFieldWithContext(StorageContextCLI, testField, testValueCLI)
	if err != nil {
		t.Fatalf("Failed to set CLI context field: %v", err)
	}

	err = SetAuthFieldWithContext(StorageContextAPI, testField, testValueProvider)
	if err != nil {
		t.Fatalf("Failed to set Provider context field: %v", err)
	}

	// Delete from CLI context
	err = DeleteAuthFieldWithContext(StorageContextCLI, testField)
	if err != nil {
		t.Fatalf("Failed to delete CLI context field: %v", err)
	}

	// Verify CLI context field is deleted
	_, err = GetAuthFieldWithContext(StorageContextCLI, testField)
	if err == nil {
		t.Errorf("CLI context field still exists after deletion")
	}

	// Verify Provider context field still exists
	valueProvider, err := GetAuthFieldWithContext(StorageContextAPI, testField)
	if err != nil {
		t.Errorf("Provider context field was deleted unexpectedly: %v", err)
	}
	if valueProvider != testValueProvider {
		t.Errorf("Provider context value changed: expected %s, got %s", testValueProvider, valueProvider)
	}

	// Cleanup
	activeProfile, _ := config.GetProfile()
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, activeProfile, testField)
}

// TestStorageContextWithProfiles tests context separation with custom profiles
func TestStorageContextWithProfiles(t *testing.T) {
	var testField authFieldKey = "test-field-profile-context"
	testProfile := makeProfileNameUnique("test-profile")

	// Make sure profile name is valid
	err := config.ValidateProfile(testProfile)
	if err != nil {
		t.Fatalf("Profile name \"%s\" is invalid: %v", testProfile, err)
	}

	testValueCLI := fmt.Sprintf("cli-value-%s", time.Now().Format(time.RFC3339))
	testValueProvider := fmt.Sprintf("provider-value-%s", time.Now().Format(time.RFC3339))

	keyring.MockInit()

	// Set values in both contexts for custom profile
	err = setAuthFieldWithProfileAndContext(StorageContextCLI, testProfile, testField, testValueCLI)
	if err != nil {
		t.Fatalf("Failed to set CLI context field for profile: %v", err)
	}

	err = setAuthFieldWithProfileAndContext(StorageContextAPI, testProfile, testField, testValueProvider)
	if err != nil {
		t.Fatalf("Failed to set Provider context field for profile: %v", err)
	}

	// Verify both contexts have correct values for the profile
	valueCLI, err := getAuthFieldWithProfileAndContext(StorageContextCLI, testProfile, testField)
	if err != nil {
		t.Fatalf("Failed to get CLI context field for profile: %v", err)
	}
	if valueCLI != testValueCLI {
		t.Errorf("CLI context value incorrect: expected %s, got %s", testValueCLI, valueCLI)
	}

	valueProvider, err := getAuthFieldWithProfileAndContext(StorageContextAPI, testProfile, testField)
	if err != nil {
		t.Fatalf("Failed to get Provider context field for profile: %v", err)
	}
	if valueProvider != testValueProvider {
		t.Errorf("Provider context value incorrect: expected %s, got %s", testValueProvider, valueProvider)
	}

	// Cleanup
	_ = deleteAuthFieldInKeyringWithContext(StorageContextCLI, testProfile, testField)
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, testProfile, testField)
	_ = deleteProfileFiles(testProfile)
}

// TestLoginLogoutWithContext tests login/logout with different contexts
func TestLoginLogoutWithContext(t *testing.T) {
	email := "test@example.com"
	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	sessionExpires := "1234567890"

	emailProvider := "provider@example.com"
	accessTokenProvider := "provider-access-token"
	refreshTokenProvider := "provider-refresh-token"
	sessionExpiresProvider := "9876543210"

	keyring.MockInit()

	// Login to CLI context
	err := LoginUserWithContext(StorageContextCLI, email, accessToken, refreshToken, sessionExpires)
	if err != nil {
		t.Fatalf("Failed to login to CLI context: %v", err)
	}

	// Login to Provider context
	err = LoginUserWithContext(StorageContextAPI, emailProvider, accessTokenProvider, refreshTokenProvider, sessionExpiresProvider)
	if err != nil {
		t.Fatalf("Failed to login to Provider context: %v", err)
	}

	// Verify CLI context credentials
	cliEmail, err := GetAuthFieldWithContext(StorageContextCLI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get CLI email: %v", err)
	}
	if cliEmail != email {
		t.Errorf("CLI email incorrect: expected %s, got %s", email, cliEmail)
	}

	cliAccessToken, err := GetAuthFieldWithContext(StorageContextCLI, ACCESS_TOKEN)
	if err != nil {
		t.Fatalf("Failed to get CLI access token: %v", err)
	}
	if cliAccessToken != accessToken {
		t.Errorf("CLI access token incorrect")
	}

	// Verify Provider context credentials
	providerEmail, err := GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get Provider email: %v", err)
	}
	if providerEmail != emailProvider {
		t.Errorf("Provider email incorrect: expected %s, got %s", emailProvider, providerEmail)
	}

	providerAccessToken, err := GetAuthFieldWithContext(StorageContextAPI, ACCESS_TOKEN)
	if err != nil {
		t.Fatalf("Failed to get Provider access token: %v", err)
	}
	if providerAccessToken != accessTokenProvider {
		t.Errorf("Provider access token incorrect")
	}

	// Logout from CLI context
	err = LogoutUserWithContext(StorageContextCLI)
	if err != nil {
		t.Fatalf("Failed to logout from CLI context: %v", err)
	}

	// Verify CLI context is logged out
	_, err = GetAuthFieldWithContext(StorageContextCLI, USER_EMAIL)
	if err == nil {
		t.Errorf("CLI context still has credentials after logout")
	}

	// Verify Provider context still has credentials
	providerEmailAfter, err := GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Provider context lost credentials after CLI logout: %v", err)
	}
	if providerEmailAfter != emailProvider {
		t.Errorf("Provider email changed after CLI logout")
	}

	// Cleanup Provider context
	err = LogoutUserWithContext(StorageContextAPI)
	if err != nil {
		t.Fatalf("Failed to logout from Provider context: %v", err)
	}
}

// TestAuthFlowWithContext tests auth flow operations with contexts
func TestAuthFlowWithContext(t *testing.T) {
	keyring.MockInit()

	// Set different auth flows for different contexts
	err := SetAuthFlowWithContext(StorageContextCLI, AUTH_FLOW_USER_TOKEN)
	if err != nil {
		t.Fatalf("Failed to set CLI auth flow: %v", err)
	}

	err = SetAuthFlowWithContext(StorageContextAPI, AUTH_FLOW_SERVICE_ACCOUNT_KEY)
	if err != nil {
		t.Fatalf("Failed to set Provider auth flow: %v", err)
	}

	// Verify CLI context auth flow
	cliFlow, err := GetAuthFlowWithContext(StorageContextCLI)
	if err != nil {
		t.Fatalf("Failed to get CLI auth flow: %v", err)
	}
	if cliFlow != AUTH_FLOW_USER_TOKEN {
		t.Errorf("CLI auth flow incorrect: expected %s, got %s", AUTH_FLOW_USER_TOKEN, cliFlow)
	}

	// Verify Provider context auth flow
	providerFlow, err := GetAuthFlowWithContext(StorageContextAPI)
	if err != nil {
		t.Fatalf("Failed to get Provider auth flow: %v", err)
	}
	if providerFlow != AUTH_FLOW_SERVICE_ACCOUNT_KEY {
		t.Errorf("Provider auth flow incorrect: expected %s, got %s", AUTH_FLOW_SERVICE_ACCOUNT_KEY, providerFlow)
	}

	// Cleanup
	activeProfile, _ := config.GetProfile()
	_ = deleteAuthFieldInKeyringWithContext(StorageContextCLI, activeProfile, authFlowType)
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, activeProfile, authFlowType)
}

// TestGetKeyringServiceName tests the keyring service name generation
func TestGetKeyringServiceName(t *testing.T) {
	tests := []struct {
		description     string
		context         StorageContext
		profile         string
		expectedService string
	}{
		{
			description:     "CLI context, default profile",
			context:         StorageContextCLI,
			profile:         config.DefaultProfileName,
			expectedService: "stackit-cli",
		},
		{
			description:     "CLI context, custom profile",
			context:         StorageContextCLI,
			profile:         "my-profile",
			expectedService: "stackit-cli/my-profile",
		},
		{
			description:     "Provider context, default profile",
			context:         StorageContextAPI,
			profile:         config.DefaultProfileName,
			expectedService: "stackit-cli-api",
		},
		{
			description:     "Provider context, custom profile",
			context:         StorageContextAPI,
			profile:         "my-profile",
			expectedService: "stackit-cli-api/my-profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			serviceName := getKeyringServiceName(tt.context, tt.profile)
			if serviceName != tt.expectedService {
				t.Errorf("Keyring service name incorrect: expected %s, got %s", tt.expectedService, serviceName)
			}
		})
	}
}

// TestGetTextFileName tests the text file name generation
func TestGetTextFileName(t *testing.T) {
	tests := []struct {
		description  string
		context      StorageContext
		expectedName string
	}{
		{
			description:  "CLI context",
			context:      StorageContextCLI,
			expectedName: "cli-auth-storage.txt",
		},
		{
			description:  "Provider context",
			context:      StorageContextAPI,
			expectedName: "cli-api-auth-storage.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			fileName := getTextFileName(tt.context)
			if fileName != tt.expectedName {
				t.Errorf("Text file name incorrect: expected %s, got %s", tt.expectedName, fileName)
			}
		})
	}
}

// TestAuthFieldMapWithContext tests bulk operations with contexts
func TestAuthFieldMapWithContext(t *testing.T) {
	testFields := map[authFieldKey]string{
		"test-field-1": fmt.Sprintf("value-1-%s", time.Now().Format(time.RFC3339)),
		"test-field-2": fmt.Sprintf("value-2-%s", time.Now().Format(time.RFC3339)),
		"test-field-3": fmt.Sprintf("value-3-%s", time.Now().Format(time.RFC3339)),
	}

	keyring.MockInit()

	// Set fields in Provider context
	err := SetAuthFieldMapWithContext(StorageContextAPI, testFields)
	if err != nil {
		t.Fatalf("Failed to set field map in Provider context: %v", err)
	}

	// Read fields from Provider context
	readFields := make(map[authFieldKey]string)
	for key := range testFields {
		readFields[key] = ""
	}
	err = GetAuthFieldMapWithContext(StorageContextAPI, readFields)
	if err != nil {
		t.Fatalf("Failed to get field map from Provider context: %v", err)
	}

	// Verify all fields match
	for key, expectedValue := range testFields {
		if readFields[key] != expectedValue {
			t.Errorf("Field %s incorrect: expected %s, got %s", key, expectedValue, readFields[key])
		}
	}

	// Verify fields don't exist in CLI context
	for key := range testFields {
		_, err := GetAuthFieldWithContext(StorageContextCLI, key)
		if err == nil {
			t.Errorf("Field %s unexpectedly exists in CLI context", key)
		}
	}

	// Cleanup
	activeProfile, _ := config.GetProfile()
	for key := range testFields {
		_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, activeProfile, key)
	}
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

// TestProviderAuthWorkflow tests the complete provider authentication workflow
func TestProviderAuthWorkflow(t *testing.T) {
	keyring.MockInit()

	email := "provider@example.com"
	accessToken := "provider-access-token"
	refreshToken := "provider-refresh-token"
	sessionExpires := fmt.Sprintf("%d", time.Now().Add(2*time.Hour).Unix())

	// Login to provider context
	err := LoginUserWithContext(StorageContextAPI, email, accessToken, refreshToken, sessionExpires)
	if err != nil {
		t.Fatalf("Failed to login to provider context: %v", err)
	}

	// Verify provider credentials exist
	providerEmail, err := GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get provider email: %v", err)
	}
	if providerEmail != email {
		t.Errorf("Provider email incorrect: expected %s, got %s", email, providerEmail)
	}

	providerAccessToken, err := GetAuthFieldWithContext(StorageContextAPI, ACCESS_TOKEN)
	if err != nil {
		t.Fatalf("Failed to get provider access token: %v", err)
	}
	if providerAccessToken != accessToken {
		t.Errorf("Provider access token incorrect")
	}

	// Verify CLI context is empty
	_, err = GetAuthFieldWithContext(StorageContextCLI, USER_EMAIL)
	if err == nil {
		t.Errorf("CLI context should be empty but has credentials")
	}

	// Set auth flow
	err = SetAuthFlowWithContext(StorageContextAPI, AUTH_FLOW_USER_TOKEN)
	if err != nil {
		t.Fatalf("Failed to set provider auth flow: %v", err)
	}

	// Verify auth flow
	providerFlow, err := GetAuthFlowWithContext(StorageContextAPI)
	if err != nil {
		t.Fatalf("Failed to get provider auth flow: %v", err)
	}
	if providerFlow != AUTH_FLOW_USER_TOKEN {
		t.Errorf("Provider auth flow incorrect: expected %s, got %s", AUTH_FLOW_USER_TOKEN, providerFlow)
	}

	// Logout from provider context
	err = LogoutUserWithContext(StorageContextAPI)
	if err != nil {
		t.Fatalf("Failed to logout from provider context: %v", err)
	}

	// Verify provider credentials are deleted
	_, err = GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err == nil {
		t.Errorf("Provider credentials still exist after logout")
	}

	// Cleanup
	activeProfile, _ := config.GetProfile()
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, activeProfile, authFlowType)
}

// TestConcurrentCLIAndProviderAuth tests that CLI and Provider can be authenticated simultaneously
func TestConcurrentCLIAndProviderAuth(t *testing.T) {
	keyring.MockInit()

	cliEmail := "cli@example.com"
	cliAccessToken := "cli-access-token"
	cliRefreshToken := "cli-refresh-token" //nolint:gosec // test credential, not a real secret
	cliSessionExpires := fmt.Sprintf("%d", time.Now().Add(2*time.Hour).Unix())

	providerEmail := "provider@example.com"
	providerAccessToken := "provider-access-token"
	providerRefreshToken := "provider-refresh-token"
	providerSessionExpires := fmt.Sprintf("%d", time.Now().Add(3*time.Hour).Unix())

	// Login to both contexts
	err := LoginUserWithContext(StorageContextCLI, cliEmail, cliAccessToken, cliRefreshToken, cliSessionExpires)
	if err != nil {
		t.Fatalf("Failed to login to CLI context: %v", err)
	}

	err = LoginUserWithContext(StorageContextAPI, providerEmail, providerAccessToken, providerRefreshToken, providerSessionExpires)
	if err != nil {
		t.Fatalf("Failed to login to Provider context: %v", err)
	}

	// Verify CLI credentials
	gotCLIEmail, err := GetAuthFieldWithContext(StorageContextCLI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get CLI email: %v", err)
	}
	if gotCLIEmail != cliEmail {
		t.Errorf("CLI email incorrect: expected %s, got %s", cliEmail, gotCLIEmail)
	}

	gotCLIAccessToken, err := GetAuthFieldWithContext(StorageContextCLI, ACCESS_TOKEN)
	if err != nil {
		t.Fatalf("Failed to get CLI access token: %v", err)
	}
	if gotCLIAccessToken != cliAccessToken {
		t.Errorf("CLI access token incorrect")
	}

	// Verify Provider credentials
	gotProviderEmail, err := GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get Provider email: %v", err)
	}
	if gotProviderEmail != providerEmail {
		t.Errorf("Provider email incorrect: expected %s, got %s", providerEmail, gotProviderEmail)
	}

	gotProviderAccessToken, err := GetAuthFieldWithContext(StorageContextAPI, ACCESS_TOKEN)
	if err != nil {
		t.Fatalf("Failed to get Provider access token: %v", err)
	}
	if gotProviderAccessToken != providerAccessToken {
		t.Errorf("Provider access token incorrect")
	}

	// Update CLI token
	newCLIAccessToken := "cli-access-token-updated"
	err = SetAuthFieldWithContext(StorageContextCLI, ACCESS_TOKEN, newCLIAccessToken)
	if err != nil {
		t.Fatalf("Failed to update CLI access token: %v", err)
	}

	// Verify CLI token was updated
	gotCLIAccessToken, err = GetAuthFieldWithContext(StorageContextCLI, ACCESS_TOKEN)
	if err != nil {
		t.Fatalf("Failed to get updated CLI access token: %v", err)
	}
	if gotCLIAccessToken != newCLIAccessToken {
		t.Errorf("CLI access token not updated: expected %s, got %s", newCLIAccessToken, gotCLIAccessToken)
	}

	// Verify Provider token unchanged
	gotProviderAccessToken, err = GetAuthFieldWithContext(StorageContextAPI, ACCESS_TOKEN)
	if err != nil {
		t.Fatalf("Failed to get Provider access token after CLI update: %v", err)
	}
	if gotProviderAccessToken != providerAccessToken {
		t.Errorf("Provider access token changed unexpectedly: expected %s, got %s", providerAccessToken, gotProviderAccessToken)
	}

	// Logout from CLI only
	err = LogoutUserWithContext(StorageContextCLI)
	if err != nil {
		t.Fatalf("Failed to logout from CLI context: %v", err)
	}

	// Verify CLI credentials are deleted
	_, err = GetAuthFieldWithContext(StorageContextCLI, USER_EMAIL)
	if err == nil {
		t.Errorf("CLI credentials still exist after logout")
	}

	// Verify Provider credentials still exist
	gotProviderEmail, err = GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Provider credentials deleted after CLI logout: %v", err)
	}
	if gotProviderEmail != providerEmail {
		t.Errorf("Provider email changed after CLI logout")
	}

	// Cleanup
	err = LogoutUserWithContext(StorageContextAPI)
	if err != nil {
		t.Fatalf("Failed to logout from provider context: %v", err)
	}
}

// TestProviderStatusReporting tests the status reporting for provider authentication
func TestProviderStatusReporting(t *testing.T) {
	keyring.MockInit()

	// Initially not authenticated
	flow, err := GetAuthFlowWithContext(StorageContextAPI)
	if err == nil && flow != "" {
		t.Errorf("Provider should not be authenticated initially, but has flow: %s", flow)
	}

	// Login
	email := "provider@example.com"
	accessToken := "provider-access-token"
	refreshToken := "provider-refresh-token"
	sessionExpires := fmt.Sprintf("%d", time.Now().Add(2*time.Hour).Unix())

	err = LoginUserWithContext(StorageContextAPI, email, accessToken, refreshToken, sessionExpires)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	err = SetAuthFlowWithContext(StorageContextAPI, AUTH_FLOW_USER_TOKEN)
	if err != nil {
		t.Fatalf("Failed to set auth flow: %v", err)
	}

	// Verify authenticated status
	flow, err = GetAuthFlowWithContext(StorageContextAPI)
	if err != nil {
		t.Fatalf("Failed to get auth flow: %v", err)
	}
	if flow != AUTH_FLOW_USER_TOKEN {
		t.Errorf("Auth flow incorrect: expected %s, got %s", AUTH_FLOW_USER_TOKEN, flow)
	}

	gotEmail, err := GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get email: %v", err)
	}
	if gotEmail != email {
		t.Errorf("Email incorrect: expected %s, got %s", email, gotEmail)
	}

	// Logout
	err = LogoutUserWithContext(StorageContextAPI)
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}

	// Verify credentials are deleted after logout
	_, err = GetAuthFieldWithContext(StorageContextAPI, USER_EMAIL)
	if err == nil {
		t.Errorf("User email should not exist after logout")
	}

	_, err = GetAuthFieldWithContext(StorageContextAPI, ACCESS_TOKEN)
	if err == nil {
		t.Errorf("Access token should not exist after logout")
	}

	// Cleanup
	activeProfile, _ := config.GetProfile()
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, activeProfile, authFlowType)
}

// TestProviderAuthWithProfiles tests provider authentication with custom profiles
func TestProviderAuthWithProfiles(t *testing.T) {
	keyring.MockInit()

	testProfile := makeProfileNameUnique("test-profile")
	err := config.ValidateProfile(testProfile)
	if err != nil {
		t.Fatalf("Profile name \"%s\" is invalid: %v", testProfile, err)
	}

	email := "provider@example.com"
	accessToken := "provider-access-token"
	refreshToken := "provider-refresh-token"
	sessionExpires := fmt.Sprintf("%d", time.Now().Add(2*time.Hour).Unix())

	// Login to provider context with custom profile
	err = setAuthFieldWithProfileAndContext(StorageContextAPI, testProfile, USER_EMAIL, email)
	if err != nil {
		t.Fatalf("Failed to set provider email for profile: %v", err)
	}

	err = setAuthFieldWithProfileAndContext(StorageContextAPI, testProfile, ACCESS_TOKEN, accessToken)
	if err != nil {
		t.Fatalf("Failed to set provider access token for profile: %v", err)
	}

	err = setAuthFieldWithProfileAndContext(StorageContextAPI, testProfile, REFRESH_TOKEN, refreshToken)
	if err != nil {
		t.Fatalf("Failed to set provider refresh token for profile: %v", err)
	}

	err = setAuthFieldWithProfileAndContext(StorageContextAPI, testProfile, SESSION_EXPIRES_AT_UNIX, sessionExpires)
	if err != nil {
		t.Fatalf("Failed to set provider session expiry for profile: %v", err)
	}

	// Verify provider credentials for custom profile
	gotEmail, err := getAuthFieldWithProfileAndContext(StorageContextAPI, testProfile, USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get provider email for profile: %v", err)
	}
	if gotEmail != email {
		t.Errorf("Provider email incorrect: expected %s, got %s", email, gotEmail)
	}

	// Verify CLI context for same profile is empty
	_, err = getAuthFieldWithProfileAndContext(StorageContextCLI, testProfile, USER_EMAIL)
	if err == nil {
		t.Errorf("CLI context for profile should be empty but has credentials")
	}

	// Cleanup
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, testProfile, USER_EMAIL)
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, testProfile, ACCESS_TOKEN)
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, testProfile, REFRESH_TOKEN)
	_ = deleteAuthFieldInKeyringWithContext(StorageContextAPI, testProfile, SESSION_EXPIRES_AT_UNIX)
	_ = deleteProfileFiles(testProfile)
}
