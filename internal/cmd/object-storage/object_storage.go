package objectstorage

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket"
	credentialsGroup "github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/enable"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "object-storage",
		Short: "Provides functionality regarding Object Storage",
		Long:  "Provides functionality regarding Object Storage.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(bucket.NewCmd())
	cmd.AddCommand(disable.NewCmd())
	cmd.AddCommand(enable.NewCmd())
	cmd.AddCommand(credentialsGroup.NewCmd())
}
