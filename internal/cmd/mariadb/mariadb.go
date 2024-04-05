package mariadb

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mariadb",
		Short: "Provides functionality for MariaDB",
		Long:  "Provides functionality for MariaDB.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(plans.NewCmd(p))
	cmd.AddCommand(credentials.NewCmd(p))
}
