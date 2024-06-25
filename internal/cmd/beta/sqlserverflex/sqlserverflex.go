package sqlserverflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/database"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sqlserverflex TEST 2",
		Short: "Provides functionality for SQLServer Flex",
		Long:  "Provides functionality for SQLServer Flex.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(database.NewCmd(p))
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(options.NewCmd(p))
	cmd.AddCommand(user.NewCmd(p))
}
