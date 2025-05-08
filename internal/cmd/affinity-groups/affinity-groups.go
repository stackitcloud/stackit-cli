package affinity_groups

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affinity-group",
		Short: "Manage server affinity groups",
		Long:  "Manage the lifecycle of server affinity groups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(
		create.NewCmd(params),
		delete.NewCmd(params),
		describe.NewCmd(params),
		list.NewCmd(params),
	)
}
