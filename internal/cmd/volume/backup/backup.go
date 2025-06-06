package backup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/backup/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/backup/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/backup/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/backup/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/backup/restore"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/backup/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Provides functionality for volume backups",
		Long:  "Provides functionality for volume backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(restore.NewCmd(params))
}
