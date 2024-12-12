package security_group

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/security-group/group"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/security-group/rule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security-group",
		Short: "Manage security groups",
		Long:  "Manage the lifecycle of security groups and rules.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(
		rule.NewCmd(p),
		group.NewCmd(p),
	)
}
