package telemetryrouter

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/access_token"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/destination"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/instance"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "telemetryrouter",
		Short: "Provides functionality for TelemetryRouter",
		Long:  "Provides functionality for TelemetryRouter.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(destination.NewCmd(params))
	cmd.AddCommand(access_token.NewCmd(params))
}
