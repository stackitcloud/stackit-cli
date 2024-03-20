package instance

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/backups"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/clone"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for PostgreSQL Flex instances",
		Long:  "Provides functionality for PostgreSQL Flex instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(update.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(clone.NewCmd())
	cmd.AddCommand(backups.NewCmd())
}
