package instance

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for LogMe instances",
		Long:  "Provides functionality for LogMe instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
