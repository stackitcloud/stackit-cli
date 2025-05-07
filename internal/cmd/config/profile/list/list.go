package list

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all CLI configuration profiles",
		Long:  "Lists all CLI configuration profiles.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List the configuration profiles`,
				"$ stackit config profile list"),
			examples.NewExample(
				`List the configuration profiles in a json format`,
				"$ stackit config profile list --output-format json"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			model := parseInput(params.Printer, cmd)

			profiles, err := config.ListProfiles()
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}

			activeProfile, err := config.GetProfile()
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}

			outputProfiles := buildOutput(profiles, activeProfile)

			return outputResult(params.Printer, model.OutputFormat, outputProfiles)
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

func buildOutput(profiles []string, activeProfile string) []profileInfo {
	var configData []profileInfo

	// Add default profile first
	configData = append(configData, profileInfo{
		Name:   config.DefaultProfileName,
		Active: activeProfile == config.DefaultProfileName,
		Email:  auth.GetProfileEmail(config.DefaultProfileName),
	})

	for _, profile := range profiles {
		configData = append(configData, profileInfo{
			Name:   profile,
			Active: profile == activeProfile,
			Email:  auth.GetProfileEmail(profile),
		})
	}

	return configData
}

func outputResult(p *print.Printer, outputFormat string, profiles []profileInfo) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(profiles, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal config list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(profiles, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal config list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("NAME", "ACTIVE", "EMAIL")
		for _, profile := range profiles {
			// Prettify the output
			email := profile.Email
			active := ""
			if profile.Email == "" {
				email = "Not authenticated"
			}
			if profile.Active {
				active = "*"
			}
			table.AddRow(profile.Name, active, email)
			table.AddSeparator()
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
