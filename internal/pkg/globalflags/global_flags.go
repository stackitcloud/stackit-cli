package globalflags

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	AsyncFlag        = "async"
	AssumeYesFlag    = "assume-yes"
	OutputFormatFlag = "output-format"
	ProjectIdFlag    = "project-id"
	VerbosityFlag    = "verbosity"

	JSONOutputFormat   = "json"
	PrettyOutputFormat = "pretty"

	DebugVerbosity   = string(print.DebugLevel)
	InfoVerbosity    = string(print.InfoLevel)
	WarningVerbosity = string(print.WarningLevel)
	ErrorVerbosity   = string(print.ErrorLevel)

	VerbosityDefault = InfoVerbosity
)

var outputFormatFlagOptions = []string{JSONOutputFormat, PrettyOutputFormat}
var verbosityFlagOptions = []string{DebugVerbosity, InfoVerbosity, WarningVerbosity, ErrorVerbosity}

type GlobalFlagModel struct {
	Async        bool
	AssumeYes    bool
	OutputFormat string
	ProjectId    string
	Verbosity    string
}

func Configure(flagSet *pflag.FlagSet) error {
	flagSet.VarP(flags.UUIDFlag(), ProjectIdFlag, "p", "Project ID")
	err := viper.BindPFlag(config.ProjectIdKey, flagSet.Lookup(ProjectIdFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", ProjectIdFlag, err)
	}

	flagSet.VarP(flags.EnumFlag(true, "", outputFormatFlagOptions...), OutputFormatFlag, "o", fmt.Sprintf("Output format, one of %q", outputFormatFlagOptions))
	err = viper.BindPFlag(config.OutputFormatKey, flagSet.Lookup(OutputFormatFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", OutputFormatFlag, err)
	}

	flagSet.Bool(AsyncFlag, false, "If set, runs the command asynchronously")
	err = viper.BindPFlag(config.AsyncKey, flagSet.Lookup(AsyncFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", AsyncFlag, err)
	}

	flagSet.BoolP(AssumeYesFlag, "y", false, "If set, skips all confirmation prompts")

	flagSet.Var(flags.EnumFlag(true, VerbosityDefault, verbosityFlagOptions...), VerbosityFlag, fmt.Sprintf("Verbosity of the CLI, one of %q", verbosityFlagOptions))
	err = viper.BindPFlag(config.VerbosityKey, flagSet.Lookup(VerbosityFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", VerbosityFlag, err)
	}

	return nil
}

func Parse(cmd *cobra.Command, p *print.Printer) *GlobalFlagModel {
	return &GlobalFlagModel{
		Async:        viper.GetBool(config.AsyncKey),
		AssumeYes:    flags.FlagToBoolValue(cmd, AssumeYesFlag, p),
		OutputFormat: viper.GetString(config.OutputFormatKey),
		ProjectId:    viper.GetString(config.ProjectIdKey),
		Verbosity:    viper.GetString(config.VerbosityKey),
	}
}
