package kubeconfig

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/kubeconfig/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/kubeconfig/login"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Provides functionality for SKE kubeconfig",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) kubeconfig.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(login.NewCmd(p))
}
