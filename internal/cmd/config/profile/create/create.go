package create

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

const (
	profileArg = "PROFILE"

	noSetFlag        = "no-set"
	fromEmptyProfile = "empty"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NoSet            bool
	FromEmptyProfile bool
	Profile          string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", profileArg),
		Short: "Creates a CLI configuration profile",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
			"Creates a CLI configuration profile based on the currently active profile and sets it as active.",
			`The profile name can be provided via the STACKIT_CLI_PROFILE environment variable or as an argument in this command.`,
			"The environment variable takes precedence over the argument.",
			"If you do not want to set the profile as active, use the --no-set flag.",
			"If you want to create the new profile with the initial default configurations, use the --empty flag.",
		),
		Args: args.SingleArg(profileArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create a new configuration profile "my-profile" with the current configuration, setting it as the active profile`,
				"$ stackit config profile create my-profile"),
			examples.NewExample(
				`Create a new configuration profile "my-profile" with a default initial configuration and don't set it as the active profile`,
				"$ stackit config profile create my-profile --empty --no-set"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			err = config.CreateProfile(p, model.Profile, !model.NoSet, model.FromEmptyProfile)
			if err != nil {
				return fmt.Errorf("create profile: %w", err)
			}

			if model.NoSet {
				p.Info("Successfully created profile %q\n", model.Profile)
				return nil
			}

			p.Info("Successfully created and set active profile to %q\n", model.Profile)

			flow, err := auth.GetAuthFlow()
			if err != nil {
				p.Debug(print.WarningLevel, "both keyring and text file storage failed to find a valid authentication flow for the active profile")
				p.Warn("The active profile %q is not authenticated, please login using the 'stackit auth login' command.\n", model.Profile)
				return nil
			}
			p.Debug(print.DebugLevel, "found valid authentication flow for active profile: %s", flow)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(noSetFlag, false, "Do not set the profile as the active profile")
	cmd.Flags().Bool(fromEmptyProfile, false, "Create the profile with the initial default configurations")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	profile := inputArgs[0]

	err := config.ValidateProfile(profile)
	if err != nil {
		return nil, err
	}

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		Profile:          profile,
		FromEmptyProfile: flags.FlagToBoolValue(p, cmd, fromEmptyProfile),
		NoSet:            flags.FlagToBoolValue(p, cmd, noSetFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}
