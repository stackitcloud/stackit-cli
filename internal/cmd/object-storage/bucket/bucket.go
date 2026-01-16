package bucket

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket",
		Short: "Provides functionality for Object Storage buckets",
		Long:  "Provides functionality for Object Storage buckets.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
