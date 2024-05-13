package print

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestOutputf(t *testing.T) {
	tests := []struct {
		description      string
		message          string
		args             []any
		verbosity        Level
		outputFormatNone bool
	}{
		{
			description: "debug verbosity",
			message:     "Test message",
			verbosity:   DebugLevel,
		},
		{
			description: "info verbosity",
			message:     "Test message",
			verbosity:   InfoLevel,
		},
		{
			description: "info verbosity - with args",
			message:     "Test message with args: %s, %s",
			args:        []any{"arg1", "arg2"},
			verbosity:   DebugLevel,
		},
		{
			description: "warning verbosity",
			message:     "Test message",
			verbosity:   WarningLevel,
		},
		{
			description: "error verbosity",
			message:     "Test message",
			verbosity:   ErrorLevel,
		},
		{
			description:      "output format none",
			message:          "Test message",
			verbosity:        InfoLevel,
			outputFormatNone: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOutput(&buf)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}
			viper.Reset()

			if tt.outputFormatNone {
				viper.Set(outputFormatKey, NoneOutputFormat)
			}

			if len(tt.args) == 0 {
				p.Outputf(tt.message)
			} else {
				p.Outputf(tt.message, tt.args...)
			}

			var expectedOutput string
			if !tt.outputFormatNone {
				expectedOutput = tt.message
				if len(tt.args) > 0 {
					expectedOutput = fmt.Sprintf(tt.message, tt.args...)
				}
			}

			output := buf.String()
			if output != expectedOutput {
				t.Errorf("unexpected output: got %q, want %q", output, expectedOutput)
			}
		})
	}
}

func TestOutputln(t *testing.T) {
	tests := []struct {
		description      string
		message          string
		verbosity        Level
		outputFormatNone bool
	}{
		{
			description: "debug verbosity",
			message:     "Test message",
			verbosity:   DebugLevel,
		},
		{
			description: "info verbosity",
			message:     "Test message",
			verbosity:   InfoLevel,
		},
		{
			description: "warning verbosity",
			message:     "Test message",
			verbosity:   WarningLevel,
		},
		{
			description: "error verbosity",
			message:     "Test message",
			verbosity:   ErrorLevel,
		},
		{
			description:      "output format none",
			message:          "Test message",
			verbosity:        InfoLevel,
			outputFormatNone: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOutput(&buf)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}
			viper.Reset()

			if tt.outputFormatNone {
				viper.Set(outputFormatKey, NoneOutputFormat)
			}

			p.Outputln(tt.message)

			var expectedOutput string
			if !tt.outputFormatNone {
				expectedOutput = fmt.Sprintf("%s\n", tt.message)
			}

			output := buf.String()
			if output != expectedOutput {
				t.Errorf("unexpected output: got %q, want %q", output, expectedOutput)
			}
		})
	}
}

func TestPagerDisplay(t *testing.T) {
	tests := []struct {
		description      string
		content          string
		verbosity        Level
		outputFormatNone bool
	}{
		{
			description: "debug verbosity",
			content:     "Test message",
			verbosity:   DebugLevel,
		},
		{
			description: "info verbosity",
			content:     "Test message",
			verbosity:   InfoLevel,
		},
		{
			description: "warning verbosity",
			content:     "Test message",
			verbosity:   WarningLevel,
		},
		{
			description: "error verbosity",
			content:     "Test message",
			verbosity:   ErrorLevel,
		},
		{
			description:      "output format none",
			content:          "Test message",
			verbosity:        InfoLevel,
			outputFormatNone: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOutput(&buf)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}
			viper.Reset()

			if tt.outputFormatNone {
				viper.Set(outputFormatKey, NoneOutputFormat)
			}

			err := p.PagerDisplay(tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			var expectedOutput string
			if !tt.outputFormatNone {
				expectedOutput = tt.content
			}

			output := buf.String()
			if output != expectedOutput {
				t.Errorf("unexpected output: got %q, want %q", output, expectedOutput)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		description string
		message     string
		args        []any
		verbosity   Level
		expectsLog  bool
		logLevel    Level
	}{
		{
			description: "debug verbosity - debug log",
			message:     "Test message",
			verbosity:   DebugLevel,
			expectsLog:  true,
			logLevel:    DebugLevel,
		},
		{
			description: "debug verbosity - info log",
			message:     "Test message",
			verbosity:   DebugLevel,
			expectsLog:  true,
			logLevel:    InfoLevel,
		},
		{
			description: "debug verbosity - warning log",
			message:     "Test message",
			verbosity:   DebugLevel,
			expectsLog:  true,
			logLevel:    WarningLevel,
		},
		{
			description: "debug verbosity - error log",
			message:     "Test message",
			verbosity:   DebugLevel,
			expectsLog:  true,
			logLevel:    ErrorLevel,
		},
		{
			description: "debug verbosity - error log with args",
			message:     "Test message",
			args:        []any{"arg1", "arg2"},
			verbosity:   DebugLevel,
			expectsLog:  true,
			logLevel:    ErrorLevel,
		},
		{
			description: "info verbosity",
			message:     "Test message",
			verbosity:   InfoLevel,
			expectsLog:  false,
			logLevel:    DebugLevel,
		},
		{
			description: "warning verbosity",
			message:     "Test message",
			verbosity:   WarningLevel,
			expectsLog:  false,
			logLevel:    DebugLevel,
		},
		{
			description: "error verbosity",
			message:     "Test message",
			verbosity:   ErrorLevel,
			expectsLog:  false,
			logLevel:    DebugLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOutput(&buf)
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
			slog.SetDefault(logger)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			if len(tt.args) == 0 {
				p.Debug(tt.logLevel, tt.message)
			}
			if len(tt.args) > 0 {
				p.Debug(tt.logLevel, tt.message, tt.args...)
			}

			output := buf.String()
			// We only check if a log is printed, as the content of the log as fields that change such as the time
			if tt.expectsLog && output == "" {
				t.Errorf("expected a log but got none")
			}
			if !tt.expectsLog && output != "" {
				t.Errorf("got log when it wasn't expected: got %q", output)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		description string
		message     string
		verbosity   Level
		shouldPrint bool
	}{
		{
			description: "debug verbosity",
			message:     "Test message",
			verbosity:   DebugLevel,
			shouldPrint: true,
		},
		{
			description: "info verbosity",
			message:     "Test message",
			verbosity:   InfoLevel,
			shouldPrint: true,
		},
		{
			description: "warning verbosity",
			message:     "Test message",
			verbosity:   WarningLevel,
			shouldPrint: false,
		},
		{
			description: "error verbosity",
			message:     "Test message",
			verbosity:   ErrorLevel,
			shouldPrint: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOutput(&buf)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			p.Info(tt.message)

			expectedOutput := tt.message
			output := buf.String()
			if tt.shouldPrint {
				if output != expectedOutput {
					t.Errorf("unexpected output: got %q, want %q", output, expectedOutput)
				}
			} else {
				if output != "" {
					t.Errorf("unexpected output: got %q, want %q", output, "")
				}
			}
		})
	}
}

func TestWarn(t *testing.T) {
	tests := []struct {
		description string
		message     string
		verbosity   Level
		shouldPrint bool
	}{
		{
			description: "debug verbosity",
			message:     "Test message",
			verbosity:   DebugLevel,
			shouldPrint: true,
		},
		{
			description: "info verbosity",
			message:     "Test message",
			verbosity:   InfoLevel,
			shouldPrint: true,
		},
		{
			description: "warning verbosity",
			message:     "Test message",
			verbosity:   WarningLevel,
			shouldPrint: true,
		},
		{
			description: "error verbosity",
			message:     "Test message",
			verbosity:   ErrorLevel,
			shouldPrint: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOutput(&buf)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			p.Warn(tt.message)

			expectedOutput := fmt.Sprintf("Warning: %s", tt.message)
			output := buf.String()
			if tt.shouldPrint {
				if output != expectedOutput {
					t.Errorf("unexpected output: got %q, want %q", output, expectedOutput)
				}
			} else {
				if output != "" {
					t.Errorf("unexpected output: got %q, want %q", output, "")
				}
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		description string
		message     string
		verbosity   Level
		shouldPrint bool
	}{
		{
			description: "debug verbosity",
			message:     "Test message",
			verbosity:   DebugLevel,
			shouldPrint: true,
		},
		{
			description: "info verbosity",
			message:     "Test message",
			verbosity:   InfoLevel,
			shouldPrint: true,
		},
		{
			description: "warning verbosity",
			message:     "Test message",
			verbosity:   WarningLevel,
			shouldPrint: true,
		},
		{
			description: "error verbosity",
			message:     "Test message",
			verbosity:   ErrorLevel,
			shouldPrint: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOutput(&buf)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			p.Error(tt.message)

			expectedOutput := fmt.Sprintf("Error: %s\n", tt.message)
			output := buf.String()
			if tt.shouldPrint {
				if output != expectedOutput {
					t.Errorf("unexpected output: got %q, want %q", output, expectedOutput)
				}
			} else {
				if output != "" {
					t.Errorf("unexpected output: got %q, want %q", output, "")
				}
			}
		})
	}
}

func TestPromptForConfirmation(t *testing.T) {
	tests := []struct {
		description string
		input       string
		verbosity   Level
		isValid     bool
		isAborted   bool
	}{
		// Note: Some of these inputs have normal spaces, others have tabs
		{
			description: "yes - simple 1",
			input:       "y\n",
			verbosity:   DebugLevel,
			isValid:     true,
		},
		{
			description: "yes - simple 2",
			input:       "  Y  \r\n",
			verbosity:   DebugLevel,
			isValid:     true,
		},
		{
			description: "yes - simple 3",
			input:       "	yes\n",
			verbosity:   DebugLevel,
			isValid:     true,
		},
		{
			description: "yes - simple 4",
			input:       "YES\n",
			verbosity:   DebugLevel,
			isValid:     true,
		},
		{
			description: "yes - retries 1",
			input:       "yrs\nyes\n",
			verbosity:   DebugLevel,
			isValid:     true,
		},
		{
			description: "yes - retries 2",
			input:       "foo\nbar  \n	y\n",
			verbosity:   DebugLevel,
			isValid:     true,
		},
		{
			description: "yes - retries 3",
			input:       "foo\r\nbar  \nY	\n",
			verbosity:   DebugLevel,
			isValid:     true,
		},
		{
			description: "no - simple 1",
			input:       "n\n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - simple 2",
			input:       "  N	\r\n",
			isValid:     false,
			verbosity:   DebugLevel,
			isAborted:   true,
		},
		{
			description: "no - simple 3",
			input:       "no\n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - simple 4",
			input:       "  \n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - simple 5",
			input:       "  \r\n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 1",
			input:       "ni\n no	\n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 2",
			input:       "foo\nbar\nn\n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 3",
			input:       "foo\r\nbar\nN\n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 4",
			input:       "m\n  \n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "no - retries 5",
			input:       "m\r\n	\r\n",
			verbosity:   DebugLevel,
			isValid:     false,
			isAborted:   true,
		},
		{
			description: "max retries 1",
			input:       "foo\nbar\nbaz\n",
			verbosity:   DebugLevel,
			isValid:     false,
		},
		{
			description: "max retries 2",
			input:       "foo\r\nbar\r\nbaz\r\n",
			verbosity:   DebugLevel,
			isValid:     false,
		},
		{
			description: "max retries 3",
			input:       "foo\nbar\nbaz\ny\n",
			verbosity:   DebugLevel,
			isValid:     false,
		},
		{
			description: "no input",
			input:       "",
			verbosity:   DebugLevel,
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
			cmd.SetErr(io.Discard)
			cmd.SetIn(buffer)

			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			err = p.PromptForConfirmation("")

			if tt.isValid && err != nil {
				t.Errorf("should not have failed: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("should have failed")
			}
			if tt.isAborted && !errors.Is(err, errAborted) {
				t.Errorf("should have returned aborted error, instead returned: %v", err)
			}
			if !tt.isAborted && errors.Is(err, errAborted) {
				t.Errorf("should not have returned aborted error")
			}
		})
	}
}

func TestIsVerbosityDebug(t *testing.T) {
	tests := []struct {
		description string
		verbosity   Level
		expected    bool
	}{
		{
			description: "debug verbosity",
			verbosity:   DebugLevel,
			expected:    true,
		},
		{
			description: "info verbosity",
			verbosity:   InfoLevel,
			expected:    false,
		},
		{
			description: "warning verbosity",
			verbosity:   WarningLevel,
			expected:    false,
		},
		{
			description: "error verbosity",
			verbosity:   ErrorLevel,
			expected:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			result := p.IsVerbosityDebug()

			if result != tt.expected {
				t.Errorf("unexpected result: got %t, want %t", result, tt.expected)
			}
		})
	}
}

func TestIsVerbosityInfo(t *testing.T) {
	tests := []struct {
		description string
		verbosity   Level
		expected    bool
	}{
		{
			description: "debug verbosity",
			verbosity:   DebugLevel,
			expected:    false,
		},
		{
			description: "info verbosity",
			verbosity:   InfoLevel,
			expected:    true,
		},
		{
			description: "warning verbosity",
			verbosity:   WarningLevel,
			expected:    false,
		},
		{
			description: "error verbosity",
			verbosity:   ErrorLevel,
			expected:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			result := p.IsVerbosityInfo()

			if result != tt.expected {
				t.Errorf("unexpected result: got %t, want %t", result, tt.expected)
			}
		})
	}
}

func TestIsVerbosityWarning(t *testing.T) {
	tests := []struct {
		description string
		verbosity   Level
		expected    bool
	}{
		{
			description: "debug verbosity",
			verbosity:   DebugLevel,
			expected:    false,
		},
		{
			description: "info verbosity",
			verbosity:   InfoLevel,
			expected:    false,
		},
		{
			description: "warning verbosity",
			verbosity:   WarningLevel,
			expected:    true,
		},
		{
			description: "error verbosity",
			verbosity:   ErrorLevel,
			expected:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			result := p.IsVerbosityWarning()

			if result != tt.expected {
				t.Errorf("unexpected result: got %t, want %t", result, tt.expected)
			}
		})
	}
}

func TestIsVerbosityError(t *testing.T) {
	tests := []struct {
		description string
		verbosity   Level
		expected    bool
	}{
		{
			description: "debug verbosity",
			verbosity:   DebugLevel,
			expected:    false,
		},
		{
			description: "info verbosity",
			verbosity:   InfoLevel,
			expected:    false,
		},
		{
			description: "warning verbosity",
			verbosity:   WarningLevel,
			expected:    false,
		},
		{
			description: "error verbosity",
			verbosity:   ErrorLevel,
			expected:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			p := &Printer{
				Cmd:       cmd,
				Verbosity: tt.verbosity,
			}

			result := p.IsVerbosityError()

			if result != tt.expected {
				t.Errorf("unexpected result: got %t, want %t", result, tt.expected)
			}
		})
	}
}
