package bucket

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	AccessKeyIdFlagName     = "access-key-id"
	SecretAccessKeyFlagName = "secret-access-key"
	EndpointFlagName        = "endpoint"
)

type ClientConfig struct {
	Region          string
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
}

type ClientFlags struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
}

func ConfigureClientFlags(flagSet *pflag.FlagSet) {
	flagSet.String(EndpointFlagName, "", "Object Storage API endpoint. Defaults to https://object.storage.<region>.onstackit.cloud. Can also be set with the "+config.ObjectStorageDataPlaneCustomEndpointKey+" config key or the STACKIT_OBJECT_STORAGE_DATA_PLANE_CUSTOM_ENDPOINT environment variable.")
	flagSet.String(AccessKeyIdFlagName, "", "Object Storage access key ID. Can also be set with the "+config.ObjectStorageAccessKeyIdKey+" config key or the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID environment variable.")
	flagSet.String(SecretAccessKeyFlagName, "", "Object Storage secret access key. Can also be set with the "+config.ObjectStorageSecretAccessKeyKey+" config key or the STACKIT_OBJECT_STORAGE_SECRET_ACCESS_KEY environment variable.")
}

func resolveFlagValue(flagValue, configKey string) string {
	if flagValue != "" {
		return flagValue
	}
	return viper.GetString(configKey)
}

func ClientConfigFromValues(clientFlags ClientFlags, region string) ClientConfig {
	return ClientConfig{
		Region:          region,
		Endpoint:        resolveFlagValue(clientFlags.Endpoint, config.ObjectStorageDataPlaneCustomEndpointKey),
		AccessKeyId:     resolveFlagValue(clientFlags.AccessKeyId, config.ObjectStorageAccessKeyIdKey),
		SecretAccessKey: resolveFlagValue(clientFlags.SecretAccessKey, config.ObjectStorageSecretAccessKeyKey),
	}
}

func ClientConfigFromReference(clientFlags ClientFlags, region string, ref Reference) ClientConfig {
	cfg := ClientConfigFromValues(clientFlags, region)
	if cfg.Endpoint == "" {
		cfg.Endpoint = ref.Endpoint
	}
	if ref.Region != "" {
		cfg.Region = ref.Region
	}
	return cfg
}

func ClientFlagsFromCommand(p *print.Printer, cmd *cobra.Command) ClientFlags {
	return ClientFlags{
		Endpoint:        flags.FlagToStringValue(p, cmd, EndpointFlagName),
		AccessKeyId:     flags.FlagToStringValue(p, cmd, AccessKeyIdFlagName),
		SecretAccessKey: flags.FlagToStringValue(p, cmd, SecretAccessKeyFlagName),
	}
}

func (c ClientConfig) ResolvedEndpoint() (string, error) {
	if c.Endpoint != "" {
		return c.Endpoint, nil
	}
	if c.Region == "" {
		return "", fmt.Errorf("no --endpoint given and no region available: use --region or --endpoint")
	}
	return defaultEndpoint(c.Region), nil
}

func defaultEndpoint(region string) string {
	return fmt.Sprintf("https://object.storage.%s.onstackit.cloud", region)
}
