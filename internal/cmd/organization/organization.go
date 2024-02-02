package organization

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/role"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization",
		Short: "Provides functionality regarding organizations",
		Long: fmt.Sprintf("%s\n%s",
			"Provides functionality regarding organizations.",
			"An active STACKIT organization is the root element of the resource hierarchy and a prerequisite to use any STACKIT Cloud Resource / Service",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(member.NewCmd())
	cmd.AddCommand(role.NewCmd())
}
