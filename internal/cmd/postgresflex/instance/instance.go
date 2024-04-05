package instance

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/clone"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for PostgreSQL Flex instances",
		Long:  "Provides functionality for PostgreSQL Flex instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(clone.NewCmd(p))
}
