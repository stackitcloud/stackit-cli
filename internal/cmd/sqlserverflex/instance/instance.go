package instance

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for SQLServer Flex instances",
		Long:  "Provides functionality for SQLServer Flex instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
