package credentialsgroup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials-group",
		Short: "Provides functionality for Object Storage credentials group",
		Long:  "Provides functionality for Object Storage credentials group.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
