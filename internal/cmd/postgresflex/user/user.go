package user

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user/list"
	resetpassword "github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user/reset-password"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for PostgreSQL Flex users",
		Long:  "Provides functionality for PostgreSQL Flex users",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(update.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(resetpassword.NewCmd())
}
