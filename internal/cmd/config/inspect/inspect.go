package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:     "inspect",
	Short:   "Inspect the current CLI configuration values",
	Long:    "Inspect the current CLI configuration values",
	Example: `$ stackit config inspect`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("read config file: %w", err)
		}

		configData := viper.AllSettings()

		configJSON, err := json.MarshalIndent(configData, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal config: %w", err)
		}
		fmt.Println(string(configJSON))

		return nil
	},
}
