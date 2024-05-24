package token

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/revoke"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Provides functionality for service account tokens",
		Long:  "Provides functionality for service account tokens.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(revoke.NewCmd(p))
}
