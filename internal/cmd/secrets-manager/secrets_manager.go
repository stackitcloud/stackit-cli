package secretsmanager

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets-manager",
		Short: "Provides functionality for Secrets Manager",
		Long:  "Provides functionality for Secrets Manager.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(user.NewCmd(p))
}
