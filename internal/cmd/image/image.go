package image

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/image/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/image/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/image/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/image/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/image/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage server images",
		Long:  "Manage the lifecycle of server images.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(
		create.NewCmd(params),
		list.NewCmd(params),
		delete.NewCmd(params),
		describe.NewCmd(params),
		update.NewCmd(params),
	)
}
