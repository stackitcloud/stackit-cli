package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/rotate"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for SKE credentials",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) credentials.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(rotate.NewCmd(p))
}
