package unset

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.UnsetProfile(p)
			if err != nil {
				return fmt.Errorf("unset profile: %w", err)
			}

			p.Info("Profile unset successfully. The default profile will be used.\n")

			flow, err := auth.GetAuthFlow()
			if err != nil {
				p.Debug(print.WarningLevel, "both keyring and text file storage failed to find a valid authentication flow for the active profile")
				p.Warn("The default profile is not authenticated, please login using the 'stackit auth login' command.\n")
				return nil
			}
			p.Debug(print.DebugLevel, "found valid authentication flow for active profile: %s", flow)

			return nil
		},
	}
	return cmd
}
