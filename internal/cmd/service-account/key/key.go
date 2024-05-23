package key

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Provides functionality for service account keys",
		Long:  "Provides functionality for service account keys.",
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
	cmd.AddCommand(update.NewCmd(p))
}
