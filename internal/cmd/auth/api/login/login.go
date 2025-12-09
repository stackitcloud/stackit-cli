package login

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
		Use:   "login",
		Short: "Logs in for the STACKIT Terraform Provider and SDK",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Logs in for the STACKIT Terraform Provider and SDK using a user account.",
			"The authentication is done via a web-based authorization flow, where the command will open a browser window in which you can login to your STACKIT account.",
			"The credentials are stored separately from the CLI authentication and will be used by the STACKIT Terraform Provider and SDK."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Login for the STACKIT Terraform Provider and SDK. This command will open a browser window where you can login to your STACKIT account`,
				"$ stackit auth api login"),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			err := auth.AuthorizeUser(params.Printer, auth.StorageContextAPI, false)
			if err != nil {
				return fmt.Errorf("authorization failed: %w", err)
			}

			params.Printer.Outputln("Successfully logged in for STACKIT Terraform Provider and SDK.\n")

			return nil
		},
	}
	return cmd
}
