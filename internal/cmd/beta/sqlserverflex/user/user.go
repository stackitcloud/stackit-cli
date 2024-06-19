package user

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/user/create"
	resetpassword "github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/user/reset-password"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for SQLServer Flex users",
		Long:  "Provides functionality for SQLServer Flex users.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(resetpassword.NewCmd(p))
}
