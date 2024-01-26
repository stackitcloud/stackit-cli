package globalflags

import (
	"fmt"

	"stackit/internal/pkg/config"
	"stackit/internal/pkg/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	AsyncFlag        = "async"
	AssumeYesFlag    = "assume-yes"
	OutputFormatFlag = "output-format"
	ProjectIdFlag    = "project-id"

	JSONOutputFormat   = "json"
	PrettyOutputFormat = "pretty"
)

var outputFormatFlagOptions = []string{JSONOutputFormat, PrettyOutputFormat}

type GlobalFlagModel struct {
	Async        bool
	AssumeYes    bool
	OutputFormat string
	ProjectId    string
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
	return nil
}

func Parse(cmd *cobra.Command) *GlobalFlagModel {
	return &GlobalFlagModel{
		Async:        viper.GetBool(config.AsyncKey),
		AssumeYes:    flags.FlagToBoolValue(cmd, AssumeYesFlag),
		OutputFormat: viper.GetString(config.OutputFormatKey),
		ProjectId:    viper.GetString(config.ProjectIdKey),
	}
}
