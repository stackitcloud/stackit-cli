package security_group

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/image/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/image/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/image/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/image/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage server images",
		Long:  "Manage the lifecycle of server images.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(
		create.NewCmd(p),
		list.NewCmd(p),
		delete.NewCmd(p),
		describe.NewCmd(p),
	)
}
