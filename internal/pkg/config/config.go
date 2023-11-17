package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Supported config keys
const (
	SessionTimeLimitKey         = "stackit_session_time_limit"
	ProjectIdKey                = "stackit_project_id"
	DNSCustomEndpointKey        = "stackit_dns_custom_endpoint"
	PostgreSQLCustomEndpointKey = "stackit_postgresql_custom_endpoint"
)

const (
	ConfigFolder        = ".stackit"
	ConfigFileName      = "cli-config"
	ConfigFileExtension = "json"
)

func InitConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configFolderPath := filepath.Join(home, ConfigFolder)
	configFilePath := filepath.Join(configFolderPath, fmt.Sprintf("%s.%s", ConfigFileName, ConfigFileExtension))

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileExtension)
	viper.AddConfigPath(configFolderPath)

	err = createFolderIfNotExists(configFolderPath)
	cobra.CheckErr(err)
	err = createFileIfNotExists(configFilePath)
	cobra.CheckErr(err)

	err = viper.ReadInConfig()
	cobra.CheckErr(err)
	setConfigDefaults()

	err = viper.WriteConfigAs(configFilePath)
	cobra.CheckErr(err)

	// Needs to be done after WriteConfigAs, otherwise it would write
	// the environment variables to the config file
	viper.AutomaticEnv()
}

func createFolderIfNotExists(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func createFileIfNotExists(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err := viper.SafeWriteConfigAs(filePath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func setConfigDefaults() {
	viper.SetDefault(SessionTimeLimitKey, "2h")
}
