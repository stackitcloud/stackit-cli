package user

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/list"
	resetpassword "github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/reset-password"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for MongoDB Flex users",
		Long:  "Provides functionality for MongoDB Flex users.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(resetpassword.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
