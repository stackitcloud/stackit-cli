package list

import (
	"fmt"
	"slices"
	"sort"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/config"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/tables"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the current CLI configuration values",
		Long:  "List the current CLI configuration values",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List your active configuration`,
				"$ stackit config list"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.ReadInConfig()
			if err != nil {
				return fmt.Errorf("read config file: %w", err)
			}

			configData := viper.AllSettings()

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
				valueString, ok := value.(string)
				if !ok || valueString == "" {
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
			err = table.Display(cmd)
			if err != nil {
				return fmt.Errorf("render table: %w", err)
			}
			return nil
		},
	}
	return cmd
}
