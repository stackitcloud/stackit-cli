package security_group

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/security_group/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/security_group/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/security_group/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/security_group/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/security_group/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security-group",
		Short: "manage security groups.",
		Long:  "manage the lifecycle of security groups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(
		create.NewCmd(p),
		delete.NewCmd(p),
		describe.NewCmd(p),
		list.NewCmd(p),
		update.NewCmd(p),
	)
}
