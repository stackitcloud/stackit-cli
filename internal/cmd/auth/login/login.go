package login

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
)

const (
	portFlag = "port"
)

type inputModel struct {
	Port *int
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logs in to the STACKIT CLI",
		Long: fmt.Sprintf("%s\n%s",
			"Logs in to the STACKIT CLI using a user account.",
			"The authentication is done via a web-based authorization flow, where the command will open a browser window in which you can login to your STACKIT account."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Login to the STACKIT CLI. This command will open a browser window where you can login to your STACKIT account`,
				"$ stackit auth login"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			err = auth.AuthorizeUser(params.Printer, auth.UserAuthConfig{
				IsReauthentication: false,
				Port:               model.Port,
			})
			if err != nil {
				return fmt.Errorf("authorization failed: %w", err)
			}

			params.Printer.Outputln("Successfully logged into STACKIT CLI.\n")

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int(portFlag, 0,
		"The port on which the callback server will listen to. By default, it tries to bind a port between 8000 and 8020.\n"+
			"When a value is specified, it will only try to use the specified port. Valid values are within the range of 8000 to 8020.",
	)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	port := flags.FlagToIntPointer(p, cmd, portFlag)
	// For the CLI client only callback URLs with localhost:[8000-8020] are valid. Additional callbacks must be enabled in the backend.
	if port != nil && (*port < 8000 || *port > 8020) {
		return nil, fmt.Errorf("port must be between 8000 and 8020")
	}

	model := inputModel{
		Port: port,
	}

	p.DebugInputModel(model)
	return &model, nil
}
