package export

import (
	"fmt"
	"path/filepath"

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

	configFileExtension = "json"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ProfileName string
	FilePath    string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("export %s", profileNameArg),
		Short: "Exports a CLI configuration profile",
		Long:  "Exports a CLI configuration profile.",
		Example: examples.Build(
			examples.NewExample(
				`Export a profile with name "PROFILE_NAME" to a file in your current directory`,
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

			err = config.ExportProfile(p, model.ProfileName, model.FilePath)
			if err != nil {
				return fmt.Errorf("could not export profile: %w", err)
			}

			p.Info("Exported profile %q to %q\n", model.ProfileName, model.FilePath)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the config to the given. If unset, writes the config to you current directory with the name of the profile. E.g. '--file-path ~/my-config.json', '--file-path ~/'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	profileName := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ProfileName:     profileName,
		FilePath:        flags.FlagToStringValue(p, cmd, filePathFlag),
	}

	// If filePath contains does not contain a file name, then add a default name
	if model.FilePath == "" {
		exportFileName := fmt.Sprintf("%s.%s", model.ProfileName, configFileExtension)
		model.FilePath = filepath.Join(model.FilePath, exportFileName)
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
