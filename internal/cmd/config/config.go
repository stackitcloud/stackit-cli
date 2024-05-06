package config

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/config/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/unset"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Provides functionality for CLI configuration options",
		Long:  fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", "Provides functionality for CLI configuration options", "The configuration is stored in a file in the user's config directory, which is OS dependent.", "Windows: %APPDATA%\\stackit", "Linux: $XDG_CONFIG_HOME/stackit or $HOME/.config/stackit", "macOS: $HOME/Library/Application Support/stackit", "The configuration file is named `cli-config.json` and is created automatically when setting a configuration option."),
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(set.NewCmd(p))
	cmd.AddCommand(unset.NewCmd(p))
}
