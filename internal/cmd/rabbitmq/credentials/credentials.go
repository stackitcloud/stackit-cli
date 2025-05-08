package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq/credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for RabbitMQ credentials",
		Long:  "Provides functionality for RabbitMQ credentials.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
