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
	IdentityProviderCustomClientIdKey = "identity_provider_custom_client_id"
	AllowedUrlDomainKey               = "allowed_url_domain"

	AuthorizationCustomEndpointKey     = "authorization_custom_endpoint"
	DNSCustomEndpointKey               = "dns_custom_endpoint"
	LoadBalancerCustomEndpointKey      = "load_balancer_custom_endpoint"
	LogMeCustomEndpointKey             = "logme_custom_endpoint"
	MariaDBCustomEndpointKey           = "mariadb_custom_endpoint"
	MongoDBFlexCustomEndpointKey       = "mongodbflex_custom_endpoint"
	ObjectStorageCustomEndpointKey     = "object_storage_custom_endpoint"
	ObservabilityCustomEndpointKey     = "observability_custom_endpoint"
	OpenSearchCustomEndpointKey        = "opensearch_custom_endpoint"
	PostgresFlexCustomEndpointKey      = "postgresflex_custom_endpoint"
	RabbitMQCustomEndpointKey          = "rabbitmq_custom_endpoint"
	RedisCustomEndpointKey             = "redis_custom_endpoint"
	ResourceManagerEndpointKey         = "resource_manager_custom_endpoint"
	SecretsManagerCustomEndpointKey    = "secrets_manager_custom_endpoint"
	ServiceAccountCustomEndpointKey    = "service_account_custom_endpoint"
	ServiceEnablementCustomEndpointKey = "service_enablement_custom_endpoint"
	ServerBackupCustomEndpointKey      = "serverbackup_custom_endpoint"
	RunCommandCustomEndpointKey        = "runcommand_custom_endpoint"
	SKECustomEndpointKey               = "ske_custom_endpoint"
	SQLServerFlexCustomEndpointKey     = "sqlserverflex_custom_endpoint"
	IaaSCustomEndpointKey              = "iaas_custom_endpoint"
	TokenCustomEndpointKey             = "token_custom_endpoint"

	ProjectNameKey     = "project_name"
	DefaultProfileName = "default"

	AsyncDefault            = false
	SessionTimeLimitDefault = "2h"

	AllowedUrlDomainDefault = "stackit.cloud"
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
	IdentityProviderCustomClientIdKey,
	AllowedUrlDomainKey,

	DNSCustomEndpointKey,
	LoadBalancerCustomEndpointKey,
	LogMeCustomEndpointKey,
	MariaDBCustomEndpointKey,
	ObjectStorageCustomEndpointKey,
	OpenSearchCustomEndpointKey,
	PostgresFlexCustomEndpointKey,
	ResourceManagerEndpointKey,
	ObservabilityCustomEndpointKey,
	AuthorizationCustomEndpointKey,
	MongoDBFlexCustomEndpointKey,
	RabbitMQCustomEndpointKey,
	RedisCustomEndpointKey,
	ResourceManagerEndpointKey,
	SecretsManagerCustomEndpointKey,
	ServiceAccountCustomEndpointKey,
	ServiceEnablementCustomEndpointKey,
	ServerBackupCustomEndpointKey,
	RunCommandCustomEndpointKey,
	SKECustomEndpointKey,
	SQLServerFlexCustomEndpointKey,
	IaaSCustomEndpointKey,
	TokenCustomEndpointKey,
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
	viper.SetDefault(IdentityProviderCustomClientIdKey, "")
	viper.SetDefault(AllowedUrlDomainKey, AllowedUrlDomainDefault)
	viper.SetDefault(DNSCustomEndpointKey, "")
	viper.SetDefault(ObservabilityCustomEndpointKey, "")
	viper.SetDefault(AuthorizationCustomEndpointKey, "")
	viper.SetDefault(MongoDBFlexCustomEndpointKey, "")
	viper.SetDefault(ObjectStorageCustomEndpointKey, "")
	viper.SetDefault(OpenSearchCustomEndpointKey, "")
	viper.SetDefault(PostgresFlexCustomEndpointKey, "")
	viper.SetDefault(ResourceManagerEndpointKey, "")
	viper.SetDefault(SecretsManagerCustomEndpointKey, "")
	viper.SetDefault(ServiceAccountCustomEndpointKey, "")
	viper.SetDefault(ServiceEnablementCustomEndpointKey, "")
	viper.SetDefault(ServerBackupCustomEndpointKey, "")
	viper.SetDefault(RunCommandCustomEndpointKey, "")
	viper.SetDefault(SKECustomEndpointKey, "")
	viper.SetDefault(SQLServerFlexCustomEndpointKey, "")
	viper.SetDefault(IaaSCustomEndpointKey, "")
	viper.SetDefault(TokenCustomEndpointKey, "")
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
