package login

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
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
		RunE: func(_ *cobra.Command, _ []string) error {
			err := auth.AuthorizeUser(p, false)
			if err != nil {
				return fmt.Errorf("authorization failed: %w", err)
			}

			p.Outputln("Successfully logged into STACKIT CLI.\n")

			return nil
		},
	}
	return cmd
}
