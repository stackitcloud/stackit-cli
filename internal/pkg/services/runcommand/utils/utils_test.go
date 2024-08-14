package utils

import (
	"testing"
)

func TestParseScriptParams(t *testing.T) {
	tests := []struct {
		description    string
		input          map[string]string
		expectedOutput map[string]string
		isValid        bool
	}{
		{
			"base-ok",
			map[string]string{"script": "ls /"},
			map[string]string{"script": "ls /"},
			true,
		},
		{
			"not-ok-nonexistant-file-specified-for-script",
			map[string]string{"script": "@{/some/file/which/does/not/exist/and/thus/fails}"},
			nil,
			false,
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
			if output["script"] != tt.expectedOutput["script"] {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput["script"], output["script"])
			}
		})
	}
}
