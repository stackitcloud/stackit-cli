package sqlserverflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/options"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sqlserverflex",
		Short: "Provides functionality for SQLServer Flex",
		Long:  "Provides functionality for SQLServer Flex.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(options.NewCmd(p))
}
