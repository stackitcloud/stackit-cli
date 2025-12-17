package postgresflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgresflex",
		Aliases: []string{"postgresqlflex"},
		Short:   "Provides functionality for PostgreSQL Flex",
		Long:    "Provides functionality for PostgreSQL Flex.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(user.NewCmd(params))
	cmd.AddCommand(options.NewCmd(params))
	cmd.AddCommand(backup.NewCmd(params))
}
