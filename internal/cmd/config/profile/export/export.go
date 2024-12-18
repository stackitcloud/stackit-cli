package export

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

const (
	profileNameArg = "PROFILE_NAME"

	filePathFlag = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ProfileName string
	ExportPath  string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("export %s", profileNameArg),
		Short: "Exports a CLI configuration profile",
		Long:  "Exports a CLI configuration profile.",
		Example: examples.Build(
			examples.NewExample(
				`Export a profile with name "PROFILE_NAME" to the current path`,
				"$ stackit config profile export PROFILE_NAME",
			),
			examples.NewExample(
				`Export a profile with name "PROFILE_NAME"" to a specific file path FILE_PATH`,
				"$ stackit config profile export PROFILE_NAME --file-path FILE_PATH",
			),
		),
		Args: args.SingleArg(profileNameArg, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			err = config.ExportProfile(p, model.ProfileName, model.ExportPath)
			if err != nil {
				return fmt.Errorf("could not export profile: %w", err)
			}

			p.Info("Exported profile %q\n", model.ProfileName)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(filePathFlag, "", "Path where the config should be saved. E.g. '--file-path ~/config.json', '--file-path ~/'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	profileName := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ProfileName:     profileName,
		ExportPath:      flags.FlagToStringValue(p, cmd, filePathFlag),
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
