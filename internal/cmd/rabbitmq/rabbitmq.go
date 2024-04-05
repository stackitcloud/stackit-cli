package rabbitmq

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rabbitmq",
		Short: "Provides functionality for RabbitMQ",
		Long:  "Provides functionality for RabbitMQ.",
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
