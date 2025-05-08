package set

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

const (
	profileArg = "PROFILE"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Profile string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("set %s", profileArg),
		Short: "Set a CLI configuration profile",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Set a CLI configuration profile as the active profile.",
			`The profile to be used can be managed via the STACKIT_CLI_PROFILE environment variable or using the "stackit config profile set PROFILE" and "stackit config profile unset" commands.`,
			"The environment variable takes precedence over what is set via the commands.",
			"When no profile is set, the default profile is used.",
		),
		Args: args.SingleArg(profileArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Set the configuration profile "my-profile" as the active profile`,
				"$ stackit config profile set my-profile"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			profileExists, err := config.ProfileExists(model.Profile)
			if err != nil {
				return fmt.Errorf("check if profile exists: %w", err)
			}
			if !profileExists {
				return &errors.SetInexistentProfile{Profile: model.Profile}
			}

			err = config.SetProfile(params.Printer, model.Profile)
			if err != nil {
				return fmt.Errorf("set profile: %w", err)
			}

			params.Printer.Info("Successfully set active profile to %q\n", model.Profile)

			flow, err := auth.GetAuthFlow()
			if err != nil {
				params.Printer.Debug(print.WarningLevel, "both keyring and text file storage failed to find a valid authentication flow for the active profile")
				params.Printer.Warn("The active profile %q is not authenticated, please login using the 'stackit auth login' command.\n", model.Profile)
				return nil
			}
			params.Printer.Debug(print.DebugLevel, "found valid authentication flow for active profile: %s", flow)

			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	profile := inputArgs[0]

	err := config.ValidateProfile(profile)
	if err != nil {
		return nil, err
	}

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Profile:         profile,
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
