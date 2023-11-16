package config

import (
	"stackit/internal/cmd/config/list"
	"stackit/internal/cmd/config/set"
	"stackit/internal/cmd/config/unset"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "CLI configuration options",
		Long:  "CLI configuration options",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(set.NewCmd())
	cmd.AddCommand(unset.NewCmd())
}
