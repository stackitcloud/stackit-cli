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
					err = deleteAuthFieldInKeyring(key)
					if err != nil {
						t.Errorf("Post-test cleanup failed: remove field \"%s\" from keyring: %v. Please remove it manually", key, err)
					}
				} else {
					err = deleteAuthFieldInEncodedTextFile(key)
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()

			for _, assignment := range tt.valueAssignments {
				err := setAuthFieldInKeyring(assignment.key, assignment.value)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", assignment.key, assignment.value, err)
				}
				// Check that this value will be checked
				if _, ok := tt.expectedValues[assignment.key]; !ok {
					t.Fatalf("Value \"%s\" set but not checked. Please add it to 'expectedValues'", assignment.key)
				}
			}

			for key, valueExpected := range tt.expectedValues {
				value, err := getAuthFieldFromKeyring(key)
				if err != nil {
					t.Errorf("Failed to get value of \"%s\": %v", key, err)
					continue
				} else if value != valueExpected {
					t.Errorf("Value of field \"%s\" is wrong: expected \"%s\", got \"%s\"", key, valueExpected, value)
				}

				err = deleteAuthFieldInKeyring(key)
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			for _, assignment := range tt.valueAssignments {
				err := setAuthFieldInEncodedTextFile(assignment.key, assignment.value)
				if err != nil {
					t.Fatalf("Failed to set \"%s\" as \"%s\": %v", assignment.key, assignment.value, err)
				}
				// Check that this value will be checked
				if _, ok := tt.expectedValues[assignment.key]; !ok {
					t.Fatalf("Value \"%s\" set but not checked. Please add it to 'expectedValues'", assignment.key)
				}
			}

			for key, valueExpected := range tt.expectedValues {
				value, err := getAuthFieldFromEncodedTextFile(key)
				if err != nil {
					t.Errorf("Failed to get value of \"%s\": %v", key, err)
					continue
				} else if value != valueExpected {
					t.Errorf("Value of field \"%s\" is wrong: expected \"%s\", got \"%s\"", key, valueExpected, value)
				}

				err = deleteAuthFieldInEncodedTextFile(key)
				if err != nil {
					t.Errorf("Post-test cleanup failed: remove field \"%s\" from text file: %v. Please remove it manually", key, err)
				}
			}
		})
	}
}

func deleteAuthFieldInKeyring(key authFieldKey) error {
	return keyring.Delete(keyringService, string(key))
}

func deleteAuthFieldInEncodedTextFile(key authFieldKey) error {
	err := createEncodedTextFile()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	textFileDir := filepath.Join(homeDir, textFileFolderName)
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
