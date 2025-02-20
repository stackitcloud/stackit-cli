package affinity_groups

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affinity-group",
		Short: "Manage server affinity groups",
		Long:  "Manage the lifecycle of server affinity groups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(
		create.NewCmd(p),
		delete.NewCmd(p),
		describe.NewCmd(p),
		list.NewCmd(p),
	)
}
