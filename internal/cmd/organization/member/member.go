package member

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member/remove"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Provides functionality regarding organization members",
		Long:  "Provides functionality regarding organization members.",
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
