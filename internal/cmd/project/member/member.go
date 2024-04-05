package member

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/member/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/member/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/member/remove"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Provides functionality regarding project members",
		Long:  "Provides functionality regarding project members.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(add.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(remove.NewCmd(p))
}
