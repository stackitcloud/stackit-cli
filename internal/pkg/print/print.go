package print

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

type Level string

const (
	DebugLevel   Level = "debug"
	InfoLevel    Level = "info"
	WarningLevel Level = "warning"
	ErrorLevel   Level = "error"

	JSONOutputFormat   = "json"
	PrettyOutputFormat = "pretty"
	NoneOutputFormat   = "none"
)

var errAborted = errors.New("operation aborted")

type Printer struct {
	Cmd       *cobra.Command
	Verbosity Level
}

// Creates a new printer, including setting up the default logger.
func NewPrinter() *Printer {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	return &Printer{}
}

// Print an output using Printf to the defined output (falling back to Stderr if not set).
// If output format is set to none, it does nothing
func (p *Printer) Outputf(msg string, args ...any) {
	outputFormat := viper.GetString(config.OutputFormatKey)
	if outputFormat == NoneOutputFormat {
		return
	}
	p.Cmd.Printf(msg, args...)
}

// Print an output using Println to the defined output (falling back to Stderr if not set).
// If output format is set to none, it does nothing
func (p *Printer) Outputln(msg string) {
	outputFormat := viper.GetString(config.OutputFormatKey)
	if outputFormat == NoneOutputFormat {
		return
	}
	p.Cmd.Println(msg)
}

// Print a Debug level log through the "slog" package.
// If the verbosity level is not Debug, it does nothing
func (p *Printer) Debug(level Level, msg string, args ...any) {
	if !p.IsVerbosityDebug() {
		return
	}
	msg = fmt.Sprintf(msg, args...)
	switch level {
	case DebugLevel:
		slog.Debug(msg)
	case InfoLevel:
		slog.Info(msg)
	case WarningLevel:
		slog.Warn(msg)
	case ErrorLevel:
		slog.Error(msg)
	}
}

// Print an Info level output to the defined Err output (falling back to Stderr if not set).
// If the verbosity level is not Debug or Info, it does nothing.
func (p *Printer) Info(msg string, args ...any) {
	if !p.IsVerbosityDebug() && !p.IsVerbosityInfo() {
		return
	}
	p.Cmd.PrintErrf(msg, args...)
}

// Print a Warn level output to the defined Err output (falling back to Stderr if not set).
// If the verbosity level is not Debug, Info, or Warn, it does nothing.
func (p *Printer) Warn(msg string, args ...any) {
	if !p.IsVerbosityDebug() && !p.IsVerbosityInfo() && !p.IsVerbosityWarning() {
		return
	}
	warning := fmt.Sprintf(msg, args...)
	p.Cmd.PrintErrf("Warning: %s", warning)
}

// Print an Error level output to the defined Err output (falling back to Stderr if not set).
func (p *Printer) Error(msg string, args ...any) {
	err := fmt.Sprintf(msg, args...)
	p.Cmd.PrintErrln(p.Cmd.ErrPrefix(), err)
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

// Shows the content in the command's stdout using the "less" command
// If output format is set to none, it does nothing
func (p *Printer) PagerDisplay(content string) error {
	outputFormat := viper.GetString(config.OutputFormatKey)
	if outputFormat == NoneOutputFormat {
		return nil
	}
	lessCmd := exec.Command("less", "-F", "-S", "-w")
	lessCmd.Stdin = strings.NewReader(content)
	lessCmd.Stdout = p.Cmd.OutOrStdout()

	err := lessCmd.Run()
	if err != nil {
		return fmt.Errorf("run less command: %w", err)
	}
	return nil
}

// Returns True if the verbosity level is set to Debug, False otherwise.
func (p *Printer) IsVerbosityDebug() bool {
	return p.Verbosity == DebugLevel
}

// Returns True if the verbosity level is set to Info, False otherwise.
func (p *Printer) IsVerbosityInfo() bool {
	return p.Verbosity == InfoLevel
}

// Returns True if the verbosity level is set to Warning, False otherwise.
func (p *Printer) IsVerbosityWarning() bool {
	return p.Verbosity == WarningLevel
}

// Returns True if the verbosity level is set to Error, False otherwise.
func (p *Printer) IsVerbosityError() bool {
	return p.Verbosity == ErrorLevel
}
