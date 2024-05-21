package config

import (
	"testing"
)

func TestReadFileIfExists(t *testing.T) {
	tests := []struct {
		description string
		filePath    string
		exists      bool
		content     string
	}{
		{
			description: "file exists",
			filePath:    "test-data/file-with-content.txt",
			exists:      true,
			content:     "my-content",
		},
		{
			description: "file does not exist",
			filePath:    "test-data/file-does-not-exist.txt",
			content:     "",
		},
		{
			description: "empty file",
			filePath:    "test-data/empty-file.txt",
			exists:      true,
			content:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			content, exists, err := readFileIfExists(tt.filePath)
			if err != nil {
				t.Errorf("read file: %v", err)
			}
			if exists != tt.exists {
				t.Errorf("expected exists to be %t but got %t", tt.exists, exists)
			}
			if content != tt.content {
				t.Errorf("expected content to be %q but got %q", tt.content, content)
			}
		})
	}
}
