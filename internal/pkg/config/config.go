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
	AsyncKey            = "async"
	OutputFormatKey     = "output_format"
	ProjectIdKey        = "project_id"
	SessionTimeLimitKey = "session_time_limit"

	DNSCustomEndpointKey            = "dns_custom_endpoint"
	MembershipCustomEndpointKey     = "membership_custom_endpoint"
	MongoDBFlexCustomEndpointKey    = "mongodbflex_custom_endpoint"
	ServiceAccountCustomEndpointKey = "service_account_custom_endpoint"
	SKECustomEndpointKey            = "ske_custom_endpoint"
	ResourceManagerEndpointKey      = "resource_manager_custom_endpoint"

	AsyncDefault            = "false"
	SessionTimeLimitDefault = "2h"
)

// Backend config keys
const (
	configFolder        = ".stackit"
	configFileName      = "cli-config"
	configFileExtension = "json"
	ProjectNameKey      = "project_name"
)

var ConfigKeys = []string{
	AsyncKey,
	OutputFormatKey,
	ProjectIdKey,
	SessionTimeLimitKey,
	DNSCustomEndpointKey,
	MembershipCustomEndpointKey,
	MongoDBFlexCustomEndpointKey,
	ServiceAccountCustomEndpointKey,
	SKECustomEndpointKey,
	ResourceManagerEndpointKey,
}

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
	viper.SetEnvPrefix("stackit")
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

// All config keys should be set to a default value so that they can be set as an environment variable
// They will not show in the config list if they are empty
func setConfigDefaults() {
	viper.SetDefault(AsyncKey, AsyncDefault)
	viper.SetDefault(OutputFormatKey, "")
	viper.SetDefault(ProjectIdKey, "")
	viper.SetDefault(SessionTimeLimitKey, SessionTimeLimitDefault)
	viper.SetDefault(DNSCustomEndpointKey, "")
	viper.SetDefault(MembershipCustomEndpointKey, "")
	viper.SetDefault(MongoDBFlexCustomEndpointKey, "")
	viper.SetDefault(ServiceAccountCustomEndpointKey, "")
	viper.SetDefault(SKECustomEndpointKey, "")
	viper.SetDefault(ResourceManagerEndpointKey, "")
}
