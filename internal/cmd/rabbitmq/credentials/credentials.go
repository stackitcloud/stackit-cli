package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for RabbitMQ credentials",
		Long:  "Provides functionality for RabbitMQ credentials.",
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
}
