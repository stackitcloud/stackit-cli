package objectstorage

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/bucket"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials"
	credentialsGroup "github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/enable"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "object-storage",
		Short: "Provides functionality for Object Storage",
		Long:  "Provides functionality for Object Storage.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(bucket.NewCmd(p))
	cmd.AddCommand(disable.NewCmd(p))
	cmd.AddCommand(enable.NewCmd(p))
	cmd.AddCommand(credentialsGroup.NewCmd(p))
	cmd.AddCommand(credentials.NewCmd(p))
}
