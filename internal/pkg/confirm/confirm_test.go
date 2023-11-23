package confirm

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/spf13/cobra"
)

func TestPromptForConfirmation(t *testing.T) {
	tests := []struct {
		description string
		input       string
		isValid     bool
		isAborted   bool
	}{
		// Note: Some of these inputs have normal spaces, others have tabs
		{
			description: "yes - simple 1",
			input:       "y\n",
			isValid:     true,
		},
		{
			description: "yes - simple 2",
			input:       "  Y  \r\n",
			isValid:     true,
		},
		{
			description: "yes - simple 3",
			input:       "	yes\n",
			isValid:     true,
		},
		{
			description: "yes - simple 4",
			input:       "YES\n",
			isValid:     true,
		},
		{
			description: "yes - retries 1",
			input:       "yrs\nyes\n",
			isValid:     true,
		},
		{
			description: "yes - retries 2",
			input:       "foo\nbar  \n	y\n",
			isValid:     true,
		},
		{
			description: "yes - retries 3",
			input:       "foo\r\nbar  \nY	\n",
			isValid:     true,
		},
		{
			description: "no - simple 1",
			input:       "n\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - simple 2",
			input:       "  N	\r\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - simple 3",
			input:       "no\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - simple 4",
			input:       "  \n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - simple 5",
			input:       "  \r\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 1",
			input:       "ni\n no	\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 2",
			input:       "foo\nbar\nn\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 3",
			input:       "foo\r\nbar\nN\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 4",
			input:       "m\n  \n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 5",
			input:       "m\r\n	\r\n",
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "max retries 1",
			input:       "foo\nbar\nbaz\n",
			isValid:     false,
		},
		{
			description: "max retries 2",
			input:       "foo\r\nbar\r\nbaz\r\n",
			isValid:     false,
		},
		{
			description: "max retries 3",
			input:       "foo\nbar\nbaz\ny\n",
			isValid:     false,
		},
		{
			description: "no input",
			input:       "",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			_, err := buffer.WriteString(tt.input)
			if err != nil {
				t.Fatalf("failed to initialize mock input: %v", err)
			}

			cmd := &cobra.Command{}
			cmd.SetOut(io.Discard) // Suppresses console prints
			cmd.SetIn(buffer)

			err = PromptForConfirmation(cmd, "")

			if tt.isValid && err != nil {
				t.Errorf("should not have failed: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("should have failed")
			}
			if tt.isAborted && !errors.Is(err, ErrAborted) {
				t.Errorf("should have returned aborted error, instead returned: %v", err)
			}
			if !tt.isAborted && errors.Is(err, ErrAborted) {
				t.Errorf("should not have returned aborted error")
			}
		})
	}
}
