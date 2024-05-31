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

			err = deleteAuthFieldProfile(tt.activeProfile)
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
			description:   "base, default profile",
			keys:          authFieldKeys,
			activeProfile: config.DefaultProfileName,
			isValid:       true,
		},
		{
			description: "missing keys, default profile",
			keys: []authFieldKey{
				ACCESS_TOKEN,
				SERVICE_ACCOUNT_EMAIL,
			},
			activeProfile: config.DefaultProfileName,
			isValid:       true,
		},
		{
			description:   "base, custom profile",
			keys:          authFieldKeys,
			activeProfile: "test-profile",
			isValid:       true,
		},
		{
			description: "missing keys, custom profile",
			keys: []authFieldKey{
				ACCESS_TOKEN,
				SERVICE_ACCOUNT_EMAIL,
			},
			activeProfile: config.DefaultProfileName,
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

			err := DeleteProfileFromKeyring(tt.activeProfile)
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

			err = deleteAuthFieldProfile(tt.activeProfile)
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

			err = deleteAuthFieldProfile(tt.activeProfile)
			if err != nil {
				t.Fatalf("Failed to remove profile: %v", err)
			}
		})
	}
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

func deleteAuthFieldProfile(activeProfile string) error {
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
