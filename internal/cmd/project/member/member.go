package member

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/member/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/member/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/member/remove"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manages project members",
		Long:  "Manages project members.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(add.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(remove.NewCmd(params))
}
