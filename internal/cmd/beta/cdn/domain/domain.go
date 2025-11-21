package domain

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCommand(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage CDN domains",
		Long:  "Manage the lifecycle of CDN domains.",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {

}
