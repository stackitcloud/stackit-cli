package project

import (
	"fmt"
	"stackit/internal/cmd/project/create"
	"stackit/internal/cmd/project/delete"
	"stackit/internal/cmd/project/describe"
	"stackit/internal/cmd/project/list"
	"stackit/internal/cmd/project/member"
	"stackit/internal/cmd/project/role"
	"stackit/internal/cmd/project/update"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Provides functionality regarding projects",
		Long: fmt.Sprintf("%s\n%s",
			"Provides functionality regarding projects.",
			"A project is a container for resources which is the service that you can purchase from STACKIT.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(update.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(member.NewCmd())
	cmd.AddCommand(role.NewCmd())
}
