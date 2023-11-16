package user

import (
	"stackit/internal/cmd/mongodbflex/user/create"
	"stackit/internal/cmd/mongodbflex/user/delete"
	"stackit/internal/cmd/mongodbflex/user/describe"
	"stackit/internal/cmd/mongodbflex/user/list"
	resetpassword "stackit/internal/cmd/mongodbflex/user/reset-password"
	"stackit/internal/cmd/mongodbflex/user/update"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for MongoDB Flex users",
		Long:  "Provides functionality for MongoDB Flex users",
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
