package organization

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/role"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization",
		Short: "Manages organizations",
		Long: fmt.Sprintf("%s\n%s",
			"Manages organizations.",
			"An active STACKIT organization is the root element of the resource hierarchy and a prerequisite to use any STACKIT Cloud Resource / Service.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(member.NewCmd(params))
	cmd.AddCommand(role.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
