package alb

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/list"
	observabilitycredentials "github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/plans"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/pool"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/quotas"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/template"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alb",
		Short: "Manages application loadbalancers",
		Long:  "Manage the lifecycle of application loadbalancers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(
		list.NewCmd(params),
		template.NewCmd(params),
		create.NewCmd(params),
		update.NewCmd(params),
		observabilitycredentials.NewCmd(params),
		describe.NewCmd(params),
		delete.NewCmd(params),
		pool.NewCmd(params),
		plans.NewCmd(params),
		quotas.NewCmd(params),
	)
}
