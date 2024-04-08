package ske

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/kubeconfig"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/options"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ske",
		Short: "Provides functionality for SKE",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE).",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(enable.NewCmd(p))
	cmd.AddCommand(kubeconfig.NewCmd())
	cmd.AddCommand(disable.NewCmd(p))
	cmd.AddCommand(cluster.NewCmd(p))
	cmd.AddCommand(credentials.NewCmd(p))
	cmd.AddCommand(options.NewCmd(p))
}
