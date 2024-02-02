package project

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/project/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/member"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/role"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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
