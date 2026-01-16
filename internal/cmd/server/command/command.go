package command

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/command/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/command/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/command/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/command/template"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "command",
		Short: "Provides functionality for Server Command",
		Long:  "Provides functionality for Server Command.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(template.NewCmd(params))
}
