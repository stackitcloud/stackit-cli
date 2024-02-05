package ske

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/options"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ske",
		Short: "Provides functionality for SKE",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE).",
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
