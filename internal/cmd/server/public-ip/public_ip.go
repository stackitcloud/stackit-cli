package publicip

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/public-ip/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/public-ip/detach"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public-ip",
		Short: "Allows attaching/detaching public IPs to servers",
		Long:  "Allows attaching/detaching public IPs to servers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(attach.NewCmd(p))
	cmd.AddCommand(detach.NewCmd(p))
}
