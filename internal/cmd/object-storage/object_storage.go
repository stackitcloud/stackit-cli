package objectstorage

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials"
	credentialsGroup "github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "object-storage",
		Short: "Provides functionality for Object Storage",
		Long:  "Provides functionality for Object Storage.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(bucket.NewCmd(params))
	cmd.AddCommand(disable.NewCmd(params))
	cmd.AddCommand(enable.NewCmd(params))
	cmd.AddCommand(credentialsGroup.NewCmd(params))
	cmd.AddCommand(credentials.NewCmd(params))
}
