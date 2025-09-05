package common

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseLabels(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expectedMap map[string]string
		expectError bool
	}{
		{
			description: "single label",
			input:       "key1=val1",
			expectedMap: map[string]string{"key1": "val1"},
			expectError: false,
		},
		{
			description: "multiple labels",
			input:       "key1=val1,key2=val2",
			expectedMap: map[string]string{"key1": "val1", "key2": "val2"},
			expectError: false,
		},
		{
			description: "empty value",
			input:       "key1=",
			expectedMap: map[string]string{"key1": ""},
			expectError: false,
		},
		{
			description: "value with equals sign",
			input:       "key1=value=with=equals",
			expectedMap: map[string]string{"key1": "value=with=equals"},
			expectError: false,
		},
		{
			description: "special case: empty string to clear labels",
			input:       "",
			expectedMap: map[string]string{}, // Should be an empty map, not nil
			expectError: false,
		},
		{
			description: "invalid format - no equals",
			input:       "key1val1",
			expectedMap: nil,
			expectError: true,
		},
		{
			description: "invalid format - empty key",
			input:       "=val1",
			expectedMap: nil,
			expectError: true,
		},
		{
			description: "mixed valid and invalid pair",
			input:       "key1=val1,key2",
			expectedMap: nil,
			expectError: true,
		},
		{
			description: "invalid format - leading comma",
			input:       ",key1=val1",
			expectedMap: nil,
			expectError: true,
		},
		{
			description: "invalid format - trailing comma",
			input:       "key1=val1,",
			expectedMap: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			parsedMap, err := ParseLabels(tt.input)

			if !tt.expectError && err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			if tt.expectError && err == nil {
				t.Fatalf("expected an error, but got nil")
			}

			if diff := cmp.Diff(tt.expectedMap, parsedMap); diff != "" {
				t.Errorf("map mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
