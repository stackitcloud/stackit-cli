package print

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestOutputf(t *testing.T) {
	tests := []struct {
		description string
		message     string
		args        []any
		verbosity   Level
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
			description: "info - with args verbosity",
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

			if len(tt.args) == 0 {
				p.Outputf(tt.message)
			} else {
				p.Outputf(tt.message, tt.args...)
			}

			expectedOutput := tt.message
			if len(tt.args) > 0 {
				expectedOutput = fmt.Sprintf(tt.message, tt.args...)
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
		description string
		message     string
		verbosity   Level
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

			p.Outputln(tt.message)

			expectedOutput := fmt.Sprintf("%s\n", tt.message)
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

			p.Debug(tt.logLevel, tt.message)

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

			expectedOutput := fmt.Sprintf("Warning: %s\n", tt.message)
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

func TestOutOrStdout(t *testing.T) {
	tests := []struct {
		description string
		writer      io.Writer
	}{
		{
			description: "os stdout",
			writer:      os.Stdout,
		},
		{
			description: "os stderr",
			writer:      os.Stderr,
		},
		{
			description: "custom bytes buffer",
			writer:      &bytes.Buffer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.SetOutput(tt.writer)
			p := &Printer{
				Cmd:       cmd,
				Verbosity: DebugLevel,
			}

			got := p.OutOrStdout()

			expected := tt.writer
			if got != expected {
				t.Errorf("unexpected output: got %v, want %v", got, expected)
			}
		})
	}
}
