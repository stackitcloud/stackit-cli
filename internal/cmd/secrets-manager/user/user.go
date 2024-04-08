package user

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/update"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for Secrets Manager users",
		Long:  "Provides functionality for Secrets Manager users.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
