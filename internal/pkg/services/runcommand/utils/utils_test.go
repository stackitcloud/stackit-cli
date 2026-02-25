package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseScriptParams(t *testing.T) {
	tests := []struct {
		description    string
		input          *map[string]string
		expectedOutput *map[string]string
		isValid        bool
	}{
		{
			description:    "base-ok",
			input:          &map[string]string{"script": "ls /"},
			expectedOutput: &map[string]string{"script": "ls /"},
			isValid:        true,
		},
		{
			description:    "nil input",
			input:          nil,
			expectedOutput: nil,
			isValid:        true,
		},
		{
			description:    "not-ok-nonexistant-file-specified-for-script",
			input:          &map[string]string{"script": "@{/some/file/which/does/not/exist/and/thus/fails}"},
			expectedOutput: nil,
			isValid:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := ParseScriptParams(tt.input)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("ParseScriptParams() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
