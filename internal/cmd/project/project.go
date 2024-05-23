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
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Provides functionality for projects",
		Long: fmt.Sprintf("%s\n%s",
			"Provides functionality for projects.",
			"A project is a container for resources which is the service that you can purchase from STACKIT.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(member.NewCmd(p))
	cmd.AddCommand(role.NewCmd(p))
}
