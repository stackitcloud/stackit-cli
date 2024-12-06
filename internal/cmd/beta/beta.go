package beta

import (
	"fmt"

	keypair "github.com/stackitcloud/stackit-cli/internal/cmd/beta/key-pair"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network"
	networkArea "github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-area"
	networkinterface "github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-interface"
	publicip "github.com/stackitcloud/stackit-cli/internal/cmd/beta/public-ip"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/security_group"
	securitygroup "github.com/stackitcloud/stackit-cli/internal/cmd/beta/security-group"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "beta",
		Short: "Contains beta STACKIT CLI commands",
		Long: fmt.Sprintf("%s\n%s",
			"Contains beta STACKIT CLI commands.",
			"The commands under this group are still in a beta state, and functionality may be incomplete or have breaking changes."),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
		Example: examples.Build(
			examples.NewExample(
				"See the currently available beta commands",
				"$ stackit beta --help"),
			examples.NewExample(
				"Execute a beta command",
				"$ stackit beta MY_COMMAND"),
		),
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(sqlserverflex.NewCmd(p))
	cmd.AddCommand(server.NewCmd(p))
	cmd.AddCommand(networkArea.NewCmd(p))
	cmd.AddCommand(network.NewCmd(p))
	cmd.AddCommand(volume.NewCmd(p))
	cmd.AddCommand(networkinterface.NewCmd(p))
	cmd.AddCommand(publicip.NewCmd(p))
	cmd.AddCommand(security_group.NewCmd(p))
	cmd.AddCommand(securitygroup.NewCmd(p))
	cmd.AddCommand(keypair.NewCmd(p))
}
