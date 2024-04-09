package logme

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logme",
		Short: "Provides functionality for LogMe",
		Long:  "Provides functionality for LogMe.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(plans.NewCmd(p))
	cmd.AddCommand(credentials.NewCmd(p))
}
