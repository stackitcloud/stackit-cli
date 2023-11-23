package globalflags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	ProjectIdFlag    = "project-id"
	OutputFormatFlag = "output-format"
	AssumeYesFlag    = "assume-yes"
)

var outputFormatFlagOptions = []string{"default", "json", "table"}

type GlobalFlagModel struct {
	ProjectId    string
	OutputFormat string
	AssumeYes    bool
}

func Configure(flagSet *pflag.FlagSet) error {
	flagSet.Var(flags.UUIDFlag(), ProjectIdFlag, "Project ID")
	err := viper.BindPFlag(config.ProjectIdKey, flagSet.Lookup(ProjectIdFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", ProjectIdFlag, err)
	}

	flagSet.Var(flags.EnumFlag(true, outputFormatFlagOptions...), OutputFormatFlag, fmt.Sprintf("Output format, one of %q", outputFormatFlagOptions))
	err = viper.BindPFlag(config.OutputFormatKey, flagSet.Lookup(OutputFormatFlag))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", OutputFormatFlag, err)
	}

	flagSet.BoolP(AssumeYesFlag, "y", false, "If set, skips all confirmation prompts")
	return nil
}

func Parse(cmd *cobra.Command) *GlobalFlagModel {
	return &GlobalFlagModel{
		ProjectId:    viper.GetString(config.ProjectIdKey),
		OutputFormat: viper.GetString(config.OutputFormatKey),
		AssumeYes:    utils.FlagToBoolValue(cmd, AssumeYesFlag),
	}
}
