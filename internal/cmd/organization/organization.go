package organization

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/member"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization/role"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization",
		Short: "Provides functionality regarding organizations",
		Long: fmt.Sprintf("%s\n%s",
			"Provides functionality regarding organizations.",
			"An active STACKIT organization is the root element of the resource hierarchy and a prerequisite to use any STACKIT Cloud Resource / Service.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(member.NewCmd(p))
	cmd.AddCommand(role.NewCmd(p))
}
