package mongodbflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mongodbflex",
		Short: "Provides functionality for MongoDB Flex",
		Long:  "Provides functionality for MongoDB Flex.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(user.NewCmd(p))
	cmd.AddCommand(options.NewCmd(p))
}
