package testutils

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Configures the given flag and binds it to the given config key.
// Should only be used in tests
func ConfigureBindUUIDFlag(cmd *cobra.Command, flag, configKey string) error {
	cmd.Flags().Var(flags.UUIDFlag(), flag, "UUID flag used for testing")
	err := viper.BindPFlag(configKey, cmd.Flags().Lookup(flag))
	if err != nil {
		return fmt.Errorf("binding --%s flag to config: %w", flag, err)
	}
	return nil
}
