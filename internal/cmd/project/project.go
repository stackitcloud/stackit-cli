package project

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manages projects",
		Long: fmt.Sprintf("%s\n%s",
			"Provides functionality for projects.",
			"A project is a container for resources which is the service that you can purchase from STACKIT.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(member.NewCmd(params))
	cmd.AddCommand(role.NewCmd(params))
}
