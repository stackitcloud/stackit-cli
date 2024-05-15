package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/cleanup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "observability-credentials",
		Short:   "Provides functionality for Load Balancer observability credentials",
		Long:    `Provides functionality for Load Balancer observability credentials. These commands can be used to store and update existing credentials, which are valid to be used for Load Balancer observability. This means, e.g. when using Argus, first of all these credentials must be created for that Argus instance (by using "stackit argus credentials create") and then can be managed for a Load Balancer by using the commands in this group.`,
		Args:    args.NoArgs,
		Aliases: []string{"credentials"},
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(add.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(cleanup.NewCmd(p))
}
