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
	SKECustomEndpointKey        = "stackit_ske_custom_endpoint"
)

const (
	configFolder        = ".stackit"
	configFileName      = "cli-config"
	configFileExtension = "json"
)

func InitConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configFolderPath := filepath.Join(home, configFolder)
	configFilePath := filepath.Join(configFolderPath, fmt.Sprintf("%s.%s", configFileName, configFileExtension))

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)
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
