package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/cleanup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "observability-credentials",
		Short:   "Provides functionality for Load Balancer observability credentials",
		Long:    `Provides functionality for Load Balancer observability credentials. These commands can be used to store and update existing credentials, which are valid to be used for Load Balancer observability. This means, e.g. when using Observability, first of all these credentials must be created for that Observability instance (by using "stackit observability credentials create") and then can be managed for a Load Balancer by using the commands in this group.`,
		Args:    args.NoArgs,
		Aliases: []string{"credentials"},
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(add.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(cleanup.NewCmd(params))
}
