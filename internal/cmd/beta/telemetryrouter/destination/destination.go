package destination

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/destination/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/destination/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/destination/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/destination/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/destination/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/telemetryrouter/destination/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destination",
		Short: "Provides functionality for TelemetryRouter destinations",
		Long:  "Provides functionality for TelemetryRouter destinations.",
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
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(generatepayload.NewCmd(params))
}
