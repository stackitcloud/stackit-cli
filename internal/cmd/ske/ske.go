package ske

import (
	"stackit/internal/cmd/ske/cluster"
	"stackit/internal/cmd/ske/credentials"
	"stackit/internal/cmd/ske/describe"
	"stackit/internal/cmd/ske/disable"
	"stackit/internal/cmd/ske/enable"
	"stackit/internal/cmd/ske/options"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ske",
		Short: "Provides functionality for SKE",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE)",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(enable.NewCmd())
	cmd.AddCommand(disable.NewCmd())
	cmd.AddCommand(cluster.NewCmd())
	cmd.AddCommand(credentials.NewCmd())
	cmd.AddCommand(options.NewCmd())
}
