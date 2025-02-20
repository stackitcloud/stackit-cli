package rule

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Provides functionality for security group rules",
		Long:  "Provides functionality for security group rules.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
