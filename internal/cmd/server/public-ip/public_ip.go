package publicip

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/public-ip/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/public-ip/detach"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public-ip",
		Short: "Allows attaching/detaching public IPs to servers",
		Long:  "Allows attaching/detaching public IPs to servers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(attach.NewCmd(params))
	cmd.AddCommand(detach.NewCmd(params))
}
