package globalflags

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
)

const (
	ProjectIdFlag GlobalFlag = "project-id"
)

type GlobalFlag string

func (f GlobalFlag) FlagName() string {
	return string(f)
}

func ConfigureFlags(flagSet *pflag.FlagSet) error {
	flagSet.Var(flags.UUIDFlag(), ProjectIdFlag.FlagName(), "Project ID")
	err := viper.BindPFlag(config.ProjectIdKey, flagSet.Lookup(ProjectIdFlag.FlagName()))
	if err != nil {
		return fmt.Errorf("bind --%s flag to config: %w", ProjectIdFlag, err)
	}
	return nil
}

func GetString(_ GlobalFlag) string {
	return viper.GetString(config.ProjectIdKey)
}
