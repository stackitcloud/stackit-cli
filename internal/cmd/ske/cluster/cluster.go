package cluster

import (
	"stackit/internal/cmd/ske/cluster/create"
	"stackit/internal/cmd/ske/cluster/delete"
	"stackit/internal/cmd/ske/cluster/describe"
	generatepayload "stackit/internal/cmd/ske/cluster/generate-payload"
	"stackit/internal/cmd/ske/cluster/list"
	"stackit/internal/cmd/ske/cluster/update"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Provides functionality for SKE cluster",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) cluster",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(generatepayload.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
