package volume

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/list"
	performanceclass "github.com/stackitcloud/stackit-cli/internal/cmd/volume/performance-class"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/resize"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Provides functionality for volumes",
		Long:  "Provides functionality for volumes.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(resize.NewCmd(params))
	cmd.AddCommand(performanceclass.NewCmd(params))
}
