package logout

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logs the user account out of the STACKIT CLI",
		Long:  "Logs the user account out of the STACKIT CLI.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Log out of the STACKIT CLI.`,
				"$ stackit auth logout"),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			err := auth.LogoutUser()
			if err != nil {
				return fmt.Errorf("log out failed: %w", err)
			}

			params.Printer.Info("Successfully logged out of the STACKIT CLI.\n")
			return nil
		},
	}
	return cmd
}
