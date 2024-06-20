package database

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/database/create"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Provides functionality for SQLServer Flex databases",
		Long:  "Provides functionality for SQLServer Flex databases.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	// cmd.AddCommand(delete.NewCmd(p))
	// cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	// cmd.AddCommand(list.NewCmd(p))
}
