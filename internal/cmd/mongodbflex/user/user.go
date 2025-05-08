package user

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/list"
	resetpassword "github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/reset-password"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for MongoDB Flex users",
		Long:  "Provides functionality for MongoDB Flex users.",
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
	cmd.AddCommand(resetpassword.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
