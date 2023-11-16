package organization

import (
	"fmt"
	"stackit/internal/cmd/organization/member"
	"stackit/internal/cmd/organization/role"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

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
