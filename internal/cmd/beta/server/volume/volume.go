package volume

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/volume/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/volume/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/volume/detach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/volume/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/volume/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Provides functionality for Server volumes",
		Long:  "Provides functionality for Server volumes.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(attach.NewCmd(p))
	cmd.AddCommand(detach.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
}