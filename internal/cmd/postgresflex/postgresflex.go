package postgresflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgresflex",
		Aliases: []string{"postgresqlflex"},
		Short:   "Provides functionality for PostgreSQL Flex",
		Long:    "Provides functionality for PostgreSQL Flex.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(user.NewCmd(p))
	cmd.AddCommand(options.NewCmd(p))
	cmd.AddCommand(backup.NewCmd(p))
}
