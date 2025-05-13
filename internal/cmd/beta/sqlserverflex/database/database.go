package database

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/database/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/database/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/database/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/database/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Provides functionality for SQLServer Flex databases",
		Long:  "Provides functionality for SQLServer Flex databases.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
