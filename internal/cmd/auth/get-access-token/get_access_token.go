package getaccesstoken

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-access-token",
		Short: "Prints a short-lived access token.",
		Long:  "Prints a short-lived access token which can be used e.g. for API calls.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Print a short-lived access token`,
				"$ stackit auth get-access-token"),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			userSessionExpired, err := auth.UserSessionExpired()
			if err != nil {
				return err
			}
			if userSessionExpired {
				return &cliErr.SessionExpiredError{}
			}

			// Try to get a valid access token, refreshing if necessary
			accessToken, err := auth.RefreshAccessToken(params.Printer)
			if err != nil {
				return err
			}

			params.Printer.Outputf("%s\n", accessToken)
			return nil
		},
	}
	return cmd
}
