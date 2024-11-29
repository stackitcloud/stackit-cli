package volume

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/list"
	performanceclass "github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/performance-class"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/resize"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Provides functionality for volumes",
		Long:  "Provides functionality for volumes.",
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
	cmd.AddCommand(resize.NewCmd(p))
	cmd.AddCommand(performanceclass.NewCmd(p))
}
