package _import

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	nameFlag   = "name"
	configFlag = "config"
	noSetFlag  = "no-set"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ProfileName string
	Config      string
	NoSet       bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Imports a CLI configuration profile",
		Long:  "Imports a CLI configuration profile.",
		Example: examples.Build(
			examples.NewExample(
				`Import a config with name "PROFILE_NAME" from file "./config.json"`,
				"$ stackit config profile --name PROFILE_NAME --config `@./config.json`",
			),
			examples.NewExample(
				`Import a config with name "PROFILE_NAME" from file "./config.json" and set not as active`,
				"$ stackit config profile --name PROFILE_NAME --config `@./config.json` --no-set",
			),
		),
		Args: args.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			err = config.ImportProfile(p, model.ProfileName, model.Config, !model.NoSet)
			if err != nil {
				return err
			}

			p.Info("Successfully imported profile %q\n", model.ProfileName)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Profile name")
	cmd.Flags().VarP(flags.ReadFromFileFlag(), configFlag, "c", "Config to be imported")
	cmd.Flags().Bool(noSetFlag, false, "Set the imported profile not as active")

	cobra.CheckErr(cmd.MarkFlagRequired(nameFlag))
	cobra.CheckErr(cmd.MarkFlagRequired(configFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		ProfileName:     flags.FlagToStringValue(p, cmd, nameFlag),
		Config:          flags.FlagToStringValue(p, cmd, configFlag),
		NoSet:           flags.FlagToBoolValue(p, cmd, noSetFlag),
	}

	if model.Config == "" {
		return nil, &errors.FlagValidationError{
			Flag:    configFlag,
			Details: "must not be empty",
		}
	}

	if model.ProfileName == "" {
		return nil, &errors.FlagValidationError{
			Flag:    nameFlag,
			Details: "must not be empty",
		}
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return model, nil
}
