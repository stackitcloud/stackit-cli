package telemetrylink

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetrylink/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetrylink/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetrylink/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetrylink/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "telemetrylink",
		Short: "Provides functionality for Telemetry Link",
		Long:  "Provides functionality for Telemetry Link.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
