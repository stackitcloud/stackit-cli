package list

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists the current CLI configuration profiles",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s",
			"Lists the current CLI configuration profiles",
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
				`List the configuration profiles`,
				"$ stackit config profile list"),
			examples.NewExample(
				`List the configuration profiles in a json format`,
				"$ stackit config profile list --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model := parseInput(p, cmd)

			profiles, err := config.ListProfiles()
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}

			activeProfile, err := config.GetProfile()
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}

			return outputResult(p, model.OutputFormat, profiles, activeProfile)
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

type profileInfo struct {
	Name   string
	Active bool
	Email  string
}

func getProfileEmail(profile string) string {
	// Get the email from the profile
	email, err := auth.GetAuthFieldWithProfile(profile, auth.USER_EMAIL)
	if err != nil {
		return ""
	}
	if email == "" {
		email, err = auth.GetAuthFieldWithProfile(profile, auth.SERVICE_ACCOUNT_EMAIL)
		if err != nil {
			return ""
		}
	}
	return email

}

func outputResult(p *print.Printer, outputFormat string, profiles []string, activeProfile string) error {
	configData := make(map[string]profileInfo)
	for _, profile := range profiles {
		configData[profile] = profileInfo{
			Name:   profile,
			Active: profile == activeProfile,
			Email:  getProfileEmail(profile),
		}
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(configData, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal config list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(configData, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal config list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		// // Sort the config options by key
		// configKeys := make([]string, 0, len(configData))
		// for k := range configData {
		// 	configKeys = append(configKeys, k)
		// }
		// sort.Strings(configKeys)

		// table := tables.NewTable()
		// if activeProfile != "" {
		// 	table.SetTitle(fmt.Sprintf("Profile: %q", activeProfile))
		// }
		// table.SetHeader("NAME", "VALUE")
		// for _, key := range configKeys {
		// 	value := configData[key]

		// 	// Convert value to string
		// 	// (Assuming value is either string or bool)
		// 	valueString, ok := value.(string)
		// 	if !ok {
		// 		valueBool, ok := value.(bool)
		// 		if !ok {
		// 			continue
		// 		}
		// 		valueString = strconv.FormatBool(valueBool)
		// 	}

		// 	// Don't show unset values
		// 	if valueString == "" {
		// 		continue
		// 	}

		// 	// Don't show unsupported (deprecated or user-inputted) configuration options
		// 	// that might be present in the config file
		// 	if !slices.Contains(config.ConfigKeys, key) {
		// 		continue
		// 	}

		// 	// Replace "_" with "-" to match the flags
		// 	key = strings.ReplaceAll(key, "_", "-")

		// 	table.AddRow(key, valueString)
		// 	table.AddSeparator()
		// }
		// err := table.Display(p)
		// if err != nil {
		// 	return fmt.Errorf("render table: %w", err)
		// }
		return nil
	}
}
