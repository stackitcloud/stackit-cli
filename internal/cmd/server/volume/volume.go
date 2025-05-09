package volume

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/volume/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/volume/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/volume/detach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/volume/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/volume/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Provides functionality for server volumes",
		Long:  "Provides functionality for server volumes.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(attach.NewCmd(params))
	cmd.AddCommand(detach.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
