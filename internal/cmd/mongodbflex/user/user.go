package user

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/list"
	resetpassword "github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/reset-password"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for MongoDB Flex users",
		Long:  "Provides functionality for MongoDB Flex users.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(resetpassword.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
