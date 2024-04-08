package serviceaccount

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/delete"
	getjwks "github.com/stackitcloud/stackit-cli/internal/cmd/service-account/get-jwks"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-account",
		Short: "Provides functionality for service accounts",
		Long:  "Provides functionality for service accounts.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(getjwks.NewCmd(p))

	cmd.AddCommand(key.NewCmd(p))
	cmd.AddCommand(token.NewCmd(p))
}
