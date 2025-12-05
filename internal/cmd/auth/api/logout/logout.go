package logout

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logs out from the STACKIT Terraform Provider and SDK",
		Long:  "Logs out from the STACKIT Terraform Provider and SDK. This does not affect CLI authentication.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Log out from the STACKIT Terraform Provider and SDK`,
				"$ stackit auth api logout"),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			err := auth.LogoutUserWithContext(auth.StorageContextAPI)
			if err != nil {
				return fmt.Errorf("log out failed: %w", err)
			}

			params.Printer.Info("Successfully logged out from STACKIT Terraform Provider and SDK.\n")
			return nil
		},
	}
	return cmd
}
