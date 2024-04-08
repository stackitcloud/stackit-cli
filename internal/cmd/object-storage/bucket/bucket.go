package bucket

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket",
		Short: "Provides functionality for Object Storage buckets",
		Long:  "Provides functionality for Object Storage buckets.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
