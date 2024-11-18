package fileutils

import (
	"os"
	"path/filepath"
	"testing"
)

const outputFilePath = "./testPayload.json"

func TestWriteToFile(t *testing.T) {
	tests := []struct {
		description string
		content     string
		outputFile  string
	}{
		{
			description: "write into file",
			content:     "Test message",
			outputFile:  outputFilePath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := WriteToFile(tt.outputFile, tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			output, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if string(output) != tt.content {
				t.Fatalf("unexpected output: got %q, want %q", output, tt.content)
			}
		})
	}
	// Cleanup
	err := os.RemoveAll(outputFilePath)
	if err != nil {
		t.Fatalf("failed cleaning test data")
	}
}

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
			content, exists, err := ReadFileIfExists(tt.filePath)
			if err != nil {
				t.Fatalf("read file: %v", err)
			}
			if exists != tt.exists {
				t.Fatalf("expected exists to be %t but got %t", tt.exists, exists)
			}
			if content != tt.content {
				t.Fatalf("expected content to be %q but got %q", tt.content, content)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tests := []struct {
		description string
		srcExists   bool
		destExists  bool
		content     string
		isValid     bool
	}{
		{
			description: "copy file",
			srcExists:   true,
			content:     "my-content",
			isValid:     true,
		},
		{
			description: "copy empty file",
			srcExists:   true,
			content:     "",
			isValid:     true,
		},
		{
			description: "copy non-existent file",
			srcExists:   false,
			content:     "",
			isValid:     false,
		},
		{
			description: "copy file to existing file",
			srcExists:   true,
			destExists:  true,
			content:     "my-content",
			isValid:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			basePath := filepath.Join(os.TempDir(), "test-data")
			src := filepath.Join(basePath, "file-with-content.txt")
			dst := filepath.Join(basePath, "file-with-content-copy.txt")

			err := os.MkdirAll(basePath, 0o750)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			if tt.srcExists {
				err := WriteToFile(src, tt.content)
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}

			if tt.destExists {
				err := WriteToFile(dst, "existing-content")
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}

			err = CopyFile(src, dst)
			if err != nil {
				if tt.isValid {
					t.Fatalf("unexpected error: %s", err.Error())
				}
				return
			}
			if !tt.isValid {
				t.Fatalf("expected error but got none")
			}

			content, exists, err := ReadFileIfExists(dst)
			if err != nil {
				t.Fatalf("read file: %v", err)
			}

			if !exists {
				t.Fatalf("expected file to exist but it does not")
			}

			if content != tt.content {
				t.Fatalf("expected content to be %q but got %q", tt.content, content)
			}

			// Cleanup
			err = os.RemoveAll(basePath)
			if err != nil {
				t.Fatalf("failed cleaning test data")
			}
		})
	}
}
