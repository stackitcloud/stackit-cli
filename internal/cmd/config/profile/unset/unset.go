package unset

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Unset the current active CLI configuration profile",
		Long: fmt.Sprintf("%s\n%s",
			"Unset the current active CLI configuration profile.",
			"When no profile is set, the default profile will be used.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Unset the currently active configuration profile. The default profile will be used.`,
				"$ stackit config profile unset"),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			err := config.UnsetProfile(params.Printer)
			if err != nil {
				return fmt.Errorf("unset profile: %w", err)
			}

			params.Printer.Info("Profile unset successfully. The default profile will be used.\n")

			flow, err := auth.GetAuthFlow()
			if err != nil {
				params.Printer.Debug(print.WarningLevel, "both keyring and text file storage failed to find a valid authentication flow for the active profile")
				params.Printer.Warn("The default profile is not authenticated, please login using the 'stackit auth login' command.\n")
				return nil
			}
			params.Printer.Debug(print.DebugLevel, "found valid authentication flow for active profile: %s", flow)

			return nil
		},
	}
	return cmd
}
