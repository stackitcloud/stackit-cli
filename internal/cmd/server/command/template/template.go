package template

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/command/template/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/command/template/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Provides functionality for Server Command Template",
		Long:  "Provides functionality for Server Command Template.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
