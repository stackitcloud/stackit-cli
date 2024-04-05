package cluster

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Provides functionality for SKE cluster",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) cluster.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(generatepayload.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
