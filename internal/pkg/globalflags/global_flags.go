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
	RegionFlag       = "region"
	VerbosityFlag    = "verbosity"

	DebugVerbosity   = string(print.DebugLevel)
	InfoVerbosity    = string(print.InfoLevel)
	WarningVerbosity = string(print.WarningLevel)
	ErrorVerbosity   = string(print.ErrorLevel)

	VerbosityDefault = InfoVerbosity
)

var outputFormatFlagOptions = []string{print.JSONOutputFormat, print.PrettyOutputFormat, print.NoneOutputFormat, print.YAMLOutputFormat}
var verbosityFlagOptions = []string{DebugVerbosity, InfoVerbosity, WarningVerbosity, ErrorVerbosity}

type GlobalFlagModel struct {
	Async        bool
	AssumeYes    bool
	OutputFormat string
	ProjectId    string
	Region       string
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

	flagSet.String(RegionFlag, "", "Target region for region-specific requests")
	err = viper.BindPFlag(config.RegionKey, flagSet.Lookup(RegionFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", RegionFlag, err)
	}

	return nil
}

func Parse(p *print.Printer, cmd *cobra.Command) *GlobalFlagModel {
	return &GlobalFlagModel{
		Async:        viper.GetBool(config.AsyncKey),
		AssumeYes:    flags.FlagToBoolValue(p, cmd, AssumeYesFlag),
		OutputFormat: viper.GetString(config.OutputFormatKey),
		ProjectId:    viper.GetString(config.ProjectIdKey),
		Region:       viper.GetString(config.RegionKey),
		Verbosity:    viper.GetString(config.VerbosityKey),
	}
}
