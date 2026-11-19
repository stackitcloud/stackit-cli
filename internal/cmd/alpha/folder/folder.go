package folder

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/folder/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/folder/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "folder",
		Short: "Manages the lifecycle of folders.",
		Long: fmt.Sprintf("%s\n%s",
			"Manages the lifecycle of folders.",
			"A folder provides a structure for organizing projects within an organization.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
