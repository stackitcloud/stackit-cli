package member

import (
	"stackit/internal/cmd/project/member/add"
	"stackit/internal/cmd/project/member/list"
	"stackit/internal/cmd/project/member/remove"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Provides functionality regarding project members",
		Long:  "Provides functionality regarding project members",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(add.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(remove.NewCmd())
}
