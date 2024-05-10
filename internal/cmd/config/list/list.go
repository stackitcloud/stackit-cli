package list

import (
	"encoding/json"
	"fmt"

	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists the current CLI configuration values",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s",
			"Lists the current CLI configuration values, based on the following sources (in order of precedence):",
			"- Environment variable",
			`  The environment variable is the name of the setting, with underscores ("_") instead of dashes ("-") and the "STACKIT" prefix.`,
			"  Example: you can set the project ID by setting the environment variable STACKIT_PROJECT_ID.",
			"- Configuration set in CLI",
			`  These are set using the "stackit config set" command`,
			`  Example: you can set the project ID by running "stackit config set --project-id xxx"`,
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List your active configuration`,
				"$ stackit config list"),
			examples.NewExample(
				`List your active configuration in a json format`,
				"$ stackit config list --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			configData := viper.AllSettings()

			model := parseInput(p, cmd)

			activeProfile, err := config.GetProfile()
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}

			return outputResult(p, model.OutputFormat, configData, activeProfile)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) *inputModel {
	globalFlags := globalflags.Parse(p, cmd)

	return &inputModel{
		GlobalFlagModel: globalFlags,
	}
}

func outputResult(p *print.Printer, outputFormat string, configData map[string]any, activeProfile string) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		if activeProfile != "" {
			configData["active_profile"] = activeProfile
		}
		details, err := json.MarshalIndent(configData, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal config list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		if activeProfile != "" {
			p.Outputf("\n ACTIVE PROFILE: %s\n", activeProfile)
		}

		// Sort the config options by key
		configKeys := make([]string, 0, len(configData))
		for k := range configData {
			configKeys = append(configKeys, k)
		}
		sort.Strings(configKeys)

		table := tables.NewTable()
		table.SetHeader("NAME", "VALUE")
		for _, key := range configKeys {
			value := configData[key]

			// Convert value to string
			// (Assuming value is either string or bool)
			valueString, ok := value.(string)
			if !ok {
				valueBool, ok := value.(bool)
				if !ok {
					continue
				}
				valueString = strconv.FormatBool(valueBool)
			}

			// Don't show unset values
			if valueString == "" {
				continue
			}

			// Don't show unsupported (deprecated or user-inputted) configuration options
			// that might be present in the config file
			if !slices.Contains(config.ConfigKeys, key) {
				continue
			}

			// Replace "_" with "-" to match the flags
			key = strings.ReplaceAll(key, "_", "-")

			table.AddRow(key, valueString)
			table.AddSeparator()
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
