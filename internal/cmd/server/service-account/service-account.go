package serviceaccount

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/service-account/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/service-account/detach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/service-account/list"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-account",
		Short: "Allows attaching/detaching service accounts to servers",
		Long:  "Allows attaching/detaching service accounts to servers",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(attach.NewCmd(params))
	cmd.AddCommand(detach.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
