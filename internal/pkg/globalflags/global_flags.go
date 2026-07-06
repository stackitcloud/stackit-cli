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
	AsyncFlag     = "async"
	AssumeYesFlag = "assume-yes"
	ProjectIdFlag = "project-id"
	RegionFlag    = "region"

	DebugVerbosity   = string(print.DebugLevel)
	InfoVerbosity    = string(print.InfoLevel)
	WarningVerbosity = string(print.WarningLevel)
	ErrorVerbosity   = string(print.ErrorLevel)

	VerbosityDefault = InfoVerbosity
)

var (
	OutputFormatFlag = flags.StringEnumFlag(
		"output-format",
		[]string{print.JSONOutputFormat, print.PrettyOutputFormat, print.NoneOutputFormat, print.YAMLOutputFormat},
		"Output format,",
		flags.StringEnumIgnoreCase[string](),
		flags.StringEnumShortHand[string]("o"),
	)
	VerbosityFlag = flags.StringEnumFlag(
		"verbosity",
		[]string{DebugVerbosity, InfoVerbosity, WarningVerbosity, ErrorVerbosity},
		"Verbosity of the CLI,",
		flags.StringEnumDefaultValue(VerbosityDefault),
		flags.StringEnumIgnoreCase[string](),
	)
)

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

	OutputFormatFlag.Register(flagSet)
	err = viper.BindPFlag(config.OutputFormatKey, flagSet.Lookup(OutputFormatFlag.Name()))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", OutputFormatFlag.Name(), err)
	}

	flagSet.Bool(AsyncFlag, false, "If set, runs the command asynchronously")
	err = viper.BindPFlag(config.AsyncKey, flagSet.Lookup(AsyncFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", AsyncFlag, err)
	}

	flagSet.BoolP(AssumeYesFlag, "y", false, "If set, skips all confirmation prompts")
	err = viper.BindPFlag(config.AssumeYesKey, flagSet.Lookup(AssumeYesFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", AssumeYesFlag, err)
	}

	VerbosityFlag.Register(flagSet)
	err = viper.BindPFlag(config.VerbosityKey, flagSet.Lookup(VerbosityFlag.Name()))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", VerbosityFlag.Name(), err)
	}

	flagSet.String(RegionFlag, "", "Target region for region-specific requests")
	err = viper.BindPFlag(config.RegionKey, flagSet.Lookup(RegionFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", RegionFlag, err)
	}

	return nil
}

func Parse(_ *print.Printer, _ *cobra.Command) *GlobalFlagModel {
	return &GlobalFlagModel{
		Async:        viper.GetBool(config.AsyncKey),
		AssumeYes:    viper.GetBool(config.AssumeYesKey),
		OutputFormat: viper.GetString(config.OutputFormatKey),
		ProjectId:    viper.GetString(config.ProjectIdKey),
		Region:       viper.GetString(config.RegionKey),
		Verbosity:    viper.GetString(config.VerbosityKey),
	}
}
