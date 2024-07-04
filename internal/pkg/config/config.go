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
	VerbosityKey        = "verbosity"

	IdentityProviderCustomEndpointKey = "identity_provider_custom_endpoint"

	ArgusCustomEndpointKey          = "argus_custom_endpoint"
	AuthorizationCustomEndpointKey  = "authorization_custom_endpoint"
	DNSCustomEndpointKey            = "dns_custom_endpoint"
	LoadBalancerCustomEndpointKey   = "load_balancer_custom_endpoint"
	LogMeCustomEndpointKey          = "logme_custom_endpoint"
	MariaDBCustomEndpointKey        = "mariadb_custom_endpoint"
	MongoDBFlexCustomEndpointKey    = "mongodbflex_custom_endpoint"
	ObjectStorageCustomEndpointKey  = "object_storage_custom_endpoint"
	OpenSearchCustomEndpointKey     = "opensearch_custom_endpoint"
	PostgresFlexCustomEndpointKey   = "postgresflex_custom_endpoint"
	RabbitMQCustomEndpointKey       = "rabbitmq_custom_endpoint"
	RedisCustomEndpointKey          = "redis_custom_endpoint"
	ResourceManagerEndpointKey      = "resource_manager_custom_endpoint"
	SecretsManagerCustomEndpointKey = "secrets_manager_custom_endpoint"
	ServiceAccountCustomEndpointKey = "service_account_custom_endpoint"
	ServerBackupCustomEndpointKey   = "serverbackup_custom_endpoint"
	SKECustomEndpointKey            = "ske_custom_endpoint"
	SQLServerFlexCustomEndpointKey  = "sqlserverflex_custom_endpoint"

	ProjectNameKey     = "project_name"
	DefaultProfileName = "default"

	AsyncDefault            = false
	SessionTimeLimitDefault = "2h"
)

const (
	configFolder = "stackit"

	configFileName      = "cli-config"
	configFileExtension = "json"

	profileRootFolder    = "profiles"
	profileFileName      = "cli-profile"
	profileFileExtension = "txt"
)

var ConfigKeys = []string{
	AsyncKey,
	OutputFormatKey,
	ProjectIdKey,
	SessionTimeLimitKey,
	VerbosityKey,

	IdentityProviderCustomEndpointKey,

	DNSCustomEndpointKey,
	LoadBalancerCustomEndpointKey,
	LogMeCustomEndpointKey,
	MariaDBCustomEndpointKey,
	ObjectStorageCustomEndpointKey,
	OpenSearchCustomEndpointKey,
	PostgresFlexCustomEndpointKey,
	ResourceManagerEndpointKey,
	ArgusCustomEndpointKey,
	AuthorizationCustomEndpointKey,
	MongoDBFlexCustomEndpointKey,
	RabbitMQCustomEndpointKey,
	RedisCustomEndpointKey,
	ResourceManagerEndpointKey,
	SecretsManagerCustomEndpointKey,
	ServiceAccountCustomEndpointKey,
	ServerBackupCustomEndpointKey,
	SKECustomEndpointKey,
	SQLServerFlexCustomEndpointKey,
}

var defaultConfigFolderPath string
var configFolderPath string
var profileFilePath string

func InitConfig() {
	defaultConfigFolderPath = getInitialConfigDir()
	profileFilePath = getInitialProfileFilePath() // Profile file path is in the default config folder

	configProfile, err := GetProfile()
	cobra.CheckErr(err)

	configFolderPath = GetProfileFolderPath(configProfile)

	configFilePath := getConfigFilePath(configFolderPath)

	// This hack is required to allow creating the config file with `viper.WriteConfig`
	// see https://github.com/spf13/viper/issues/851#issuecomment-789393451
	viper.SetConfigFile(configFilePath)

	f, err := os.Open(configFilePath)
	if !os.IsNotExist(err) {
		if err := viper.ReadConfig(f); err != nil {
			cobra.CheckErr(err)
		}
	}
	defer func() {
		if f != nil {
			if err := f.Close(); err != nil {
				cobra.CheckErr(err)
			}
		}
	}()

	setConfigDefaults()

	viper.AutomaticEnv()
	viper.SetEnvPrefix("stackit")
}

// Write saves the config file (wrapping `viper.WriteConfig`) and ensures that its directory exists
func Write() error {
	err := os.MkdirAll(configFolderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	return viper.WriteConfig()
}

// All config keys should be set to a default value so that they can be set as an environment variable
// They will not show in the config list if they are empty
func setConfigDefaults() {
	viper.SetDefault(AsyncKey, AsyncDefault)
	viper.SetDefault(OutputFormatKey, "")
	viper.SetDefault(ProjectIdKey, "")
	viper.SetDefault(SessionTimeLimitKey, SessionTimeLimitDefault)
	viper.SetDefault(IdentityProviderCustomEndpointKey, "")
	viper.SetDefault(DNSCustomEndpointKey, "")
	viper.SetDefault(ArgusCustomEndpointKey, "")
	viper.SetDefault(AuthorizationCustomEndpointKey, "")
	viper.SetDefault(MongoDBFlexCustomEndpointKey, "")
	viper.SetDefault(ObjectStorageCustomEndpointKey, "")
	viper.SetDefault(OpenSearchCustomEndpointKey, "")
	viper.SetDefault(PostgresFlexCustomEndpointKey, "")
	viper.SetDefault(ResourceManagerEndpointKey, "")
	viper.SetDefault(SecretsManagerCustomEndpointKey, "")
	viper.SetDefault(ServiceAccountCustomEndpointKey, "")
	viper.SetDefault(ServerBackupCustomEndpointKey, "")
	viper.SetDefault(SKECustomEndpointKey, "")
	viper.SetDefault(SQLServerFlexCustomEndpointKey, "")
}

func getConfigFilePath(configFolder string) string {
	return filepath.Join(configFolder, fmt.Sprintf("%s.%s", configFileName, configFileExtension))
}

func getInitialConfigDir() string {
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	return filepath.Join(configDir, configFolder)
}

func getInitialProfileFilePath() string {
	configFolderPath := defaultConfigFolderPath
	if configFolderPath == "" {
		configFolderPath = getInitialConfigDir()
	}
	return filepath.Join(configFolderPath, fmt.Sprintf("%s.%s", profileFileName, profileFileExtension))
}
