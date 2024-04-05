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
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for PostgreSQL Flex users",
		Long:  "Provides functionality for PostgreSQL Flex users.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(resetpassword.NewCmd(p))
}
