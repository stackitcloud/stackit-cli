package print

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type VerbosityLevel string

const (
	DebugVerbosity   VerbosityLevel = "debug"
	InfoVerbosity    VerbosityLevel = "info"
	WarningVerbosity VerbosityLevel = "warning"
	ErrorVerbosity   VerbosityLevel = "error"
)

type Printer struct {
	Cmd       *cobra.Command
	Verbosity VerbosityLevel
}

// Creates a new printer, including setting up the default logger.
func NewPrinter() Printer {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	return Printer{}
}

// Print an output using Printf to the defined output (falling back to Stderr if not set).
func (p *Printer) Outputf(message string, args ...any) {
	p.Cmd.Printf(message, args...)
}

// Print an output using Println to the defined output (falling back to Stderr if not set).
func (p *Printer) Outputln(message string) {
	p.Cmd.Println(message)
}

// Print a Debug level log through the "slog" package.
// If the verbosity level is not Debug, it does nothing
func (p *Printer) Debug(message string, args ...any) {
	if p.Verbosity != DebugVerbosity {
		return
	}
	slog.Debug(message, args...)
}

// Print an Info level log to the defined Err output (falling back to Stderr if not set).
// If the verbosity level is not Debug or Info, it does nothing.
func (p *Printer) Info(message string, args ...any) {
	if p.Verbosity != DebugVerbosity && p.Verbosity != InfoVerbosity {
		return
	}
	p.Cmd.PrintErrf(message, args...)
}

// Print an Warning level log to the defined Err output (falling back to Stderr if not set).
// If the verbosity level is not Debug, Info, or Warning, it does nothing.
func (p *Printer) Warning(message string) {
	if p.Verbosity != DebugVerbosity && p.Verbosity != InfoVerbosity && p.Verbosity != WarningVerbosity {
		return
	}
	p.Cmd.PrintErrf("Warning: %s\n", message)
}

// Print an Error level log to the defined Err output (falling back to Stderr if not set).
func (p *Printer) Error(message string) {
	p.Cmd.PrintErrln(p.Cmd.ErrPrefix(), message)
}

// Returns the printer's command defined output
func (p *Printer) OutOrStdout() io.Writer {
	return p.Cmd.OutOrStdout()
}
