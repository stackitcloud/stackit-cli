package globalflags

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
)

const (
	ProjectIdFlag     GlobalFlag = "project-id"
	OutputFormatFlag  GlobalFlag = "output-format"
	JSONOutputFormat             = "json"
	TableOutputFormat            = "table"
)

var outputFormatFlagOptions = []string{JSONOutputFormat, TableOutputFormat}

type GlobalFlag string

type Model struct {
	ProjectId    string
	OutputFormat string
}

func (f GlobalFlag) FlagName() string {
	return string(f)
}

func Configure(flagSet *pflag.FlagSet) error {
	flagSet.Var(flags.UUIDFlag(), ProjectIdFlag.FlagName(), "Project ID")
	err := viper.BindPFlag(config.ProjectIdKey, flagSet.Lookup(ProjectIdFlag.FlagName()))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", ProjectIdFlag, err)
	}

	flagSet.Var(flags.EnumFlag(true, outputFormatFlagOptions...), OutputFormatFlag.FlagName(), fmt.Sprintf("Output format, one of %q", outputFormatFlagOptions))
	err = viper.BindPFlag(config.OutputFormatKey, flagSet.Lookup(OutputFormatFlag.FlagName()))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", OutputFormatFlag, err)
	}
	return nil
}

func Parse() *Model {
	return &Model{
		ProjectId:    viper.GetString(config.ProjectIdKey),
		OutputFormat: viper.GetString(config.OutputFormatKey),
	}
}
