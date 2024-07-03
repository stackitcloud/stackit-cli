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
		Long:  "Logs in to the STACKIT CLI using a user account.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Login to the STACKIT CLI. This command will open a browser window where you can login to your STACKIT account`,
				"$ stackit auth login"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			p.Warn(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n\n",
				"Starting on July 9 2024, the new STACKIT Identity Provider (IDP) will be available.",
				"On this date, we will release a new version of the STACKIT CLI that will use the new IDP for user authentication.",
				"This also means that the user authentication on STACKIT CLI versions released before July 9 2024 is no longer guaranteed to work for all services.",
				"Please make sure to update your STACKIT CLI to the latest version after July 9 2024 to ensure that you can continue to use all STACKIT services.",
				"You can find more information regarding the new IDP at https://docs.stackit.cloud/stackit/en/release-notes-23101442.html#ReleaseNotes-2024-06-21-identity-provider.",
			))

			err := auth.AuthorizeUser(p, false)
			if err != nil {
				return fmt.Errorf("authorization failed: %w", err)
			}

			p.Info("Successfully logged into STACKIT CLI.\n")
			return nil
		},
	}
	return cmd
}
