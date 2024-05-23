package fileutils

import (
	"os"
	"testing"
)

const outputFilePath = "./testPayload.json"

func TestFileOutput(t *testing.T) {
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
			err := FileOutput(tt.outputFile, tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			output, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if string(output) != tt.content {
				t.Errorf("unexpected output: got %q, want %q", output, tt.content)
			}
		})
	}
	// Cleanup
	err := os.RemoveAll(outputFilePath)
	if err != nil {
		t.Errorf("failed cleaning test data")
	}
}
