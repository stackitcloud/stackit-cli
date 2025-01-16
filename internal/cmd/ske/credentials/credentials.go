package credentials

import (
	completerotation "github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/complete-rotation"
	startrotation "github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/start-rotation"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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
	cmd.AddCommand(startrotation.NewCmd(p))
	cmd.AddCommand(completerotation.NewCmd(p))
}
