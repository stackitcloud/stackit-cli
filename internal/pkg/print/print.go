package print

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Level string

const (
	DebugLevel   Level = "debug"
	InfoLevel    Level = "info"
	WarningLevel Level = "warning"
	ErrorLevel   Level = "error"
)

var errAborted = errors.New("operation aborted")

type Printer struct {
	Cmd       *cobra.Command
	Verbosity Level
}

// Creates a new printer, including setting up the default logger.
func NewPrinter() *Printer {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	return &Printer{}
}

// Print an output using Printf to the defined output (falling back to Stderr if not set).
func (p *Printer) Outputf(msg string, args ...any) {
	p.Cmd.Printf(msg, args...)
}

// Print an output using Println to the defined output (falling back to Stderr if not set).
func (p *Printer) Outputln(msg string) {
	p.Cmd.Println(msg)
}

// Print a Debug level log through the "slog" package.
// If the verbosity level is not Debug, it does nothing
func (p *Printer) Debug(level Level, msg string, args ...any) {
	if p.Verbosity != DebugLevel {
		return
	}
	switch level {
	case DebugLevel:
		slog.Debug(msg, args...)
	case InfoLevel:
		slog.Info(msg, args...)
	case WarningLevel:
		slog.Warn(msg, args...)
	case ErrorLevel:
		slog.Error(msg, args...)
	}
}

// Print an Info level output to the defined Err output (falling back to Stderr if not set).
// If the verbosity level is not Debug or Info, it does nothing.
func (p *Printer) Info(msg string, args ...any) {
	if p.Verbosity != DebugLevel && p.Verbosity != InfoLevel {
		return
	}
	p.Cmd.PrintErrf(msg, args...)
}

// Print a Warn level output to the defined Err output (falling back to Stderr if not set).
// If the verbosity level is not Debug, Info, or Warn, it does nothing.
func (p *Printer) Warn(msg string) {
	if p.Verbosity != DebugLevel && p.Verbosity != InfoLevel && p.Verbosity != WarningLevel {
		return
	}
	p.Cmd.PrintErrf("Warning: %s\n", msg)
}

// Print an Error level output to the defined Err output (falling back to Stderr if not set).
func (p *Printer) Error(msg string) {
	p.Cmd.PrintErrln(p.Cmd.ErrPrefix(), msg)
}

// Prompts the user for confirmation.
//
// Returns nil only if the user (explicitly) answers positive.
// Returns ErrAborted if the user answers negative.
func (p *Printer) PromptForConfirmation(prompt string) error {
	question := fmt.Sprintf("%s [y/N] ", prompt)
	reader := bufio.NewReader(p.Cmd.InOrStdin())
	for i := 0; i < 3; i++ {
		p.Cmd.PrintErr(question)
		answer, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer == "y" || answer == "yes" {
			return nil
		}
		if answer == "" || answer == "n" || answer == "no" {
			return errAborted
		}
	}
	return fmt.Errorf("max number of wrong inputs")
}
