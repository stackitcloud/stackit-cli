package credentials

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	add "github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "observability-credentials",
		Short: "Provides functionality for application loadbalancer credentials",
		Long:  "Provides functionality for application loadbalancer credentials",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(add.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
