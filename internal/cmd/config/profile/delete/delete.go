package delete

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
		Use:   fmt.Sprintf("delete %s", profileArg),
		Short: "Delete a CLI configuration profile",
		Long: fmt.Sprintf("%s\n%s",
			"Delete a CLI configuration profile.",
			"If the deleted profile is the active profile, the default profile will be set to active.",
		),
		Args: args.SingleArg(profileArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete the configuration profile "my-profile"`,
				"$ stackit config profile delete my-profile"),
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
				return &errors.DeleteInexistentProfile{Profile: model.Profile}
			}

			if model.Profile == config.DefaultProfileName {
				return &errors.DeleteDefaultProfile{DefaultProfile: config.DefaultProfileName}
			}

			activeProfile, err := config.GetProfile()
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}
			if activeProfile == model.Profile {
				params.Printer.Warn("The profile you are trying to delete is the active profile. The default profile will be set to active.\n")
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete profile %q? (This cannot be undone)", model.Profile)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			err = config.DeleteProfile(params.Printer, model.Profile)
			if err != nil {
				return fmt.Errorf("delete profile: %w", err)
			}

			err = auth.DeleteProfileAuth(model.Profile)
			if err != nil {
				return fmt.Errorf("delete profile authentication: %w", err)
			}

			params.Printer.Info("Successfully deleted profile %q\n", model.Profile)

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
