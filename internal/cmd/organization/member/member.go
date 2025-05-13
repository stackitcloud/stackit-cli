package member

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member/remove"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manages organization members",
		Long:  "Manages organization members.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(add.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(remove.NewCmd(params))
}
