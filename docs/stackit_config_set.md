## stackit config set

Sets CLI configuration options

### Synopsis

Sets CLI configuration options.
All of the configuration options can be set using an environment variable, which takes precedence over what is configured using this command.
The environment variable is the name of the flag, with underscores ("_") instead of dashes ("-") and the "STACKIT" prefix.
Example: to set the project ID you can set the environment variable STACKIT_PROJECT_ID.

```
stackit config set [flags]
```

### Examples

```
  Set a project ID in your active configuration. This project ID will be used by every command (unless overridden by the "STACKIT_PROJECT_ID" environment variable)
  $ stackit config set --project-id xxx

  Set the session time limit to 1 hour
  $ stackit config set --session-time-limit 1h

  Set the DNS custom endpoint. This endpoint will be used on all calls to the DNS API (unless overridden by the "STACKIT_DNS_CUSTOM_ENDPOINT" environment variable)
  $ stackit config set --dns-custom-endpoint https://dns.stackit.cloud
```

### Options

```
      --allowed-url-domain string                                  Domain name, used for the verification of the URLs that are given in the custom identity provider endpoint and "STACKIT curl" command
      --authorization-custom-endpoint string                       Authorization API base URL, used in calls to this API
      --dns-custom-endpoint string                                 DNS API base URL, used in calls to this API
      --edge-custom-endpoint string                                Edge API base URL, used in calls to this API
  -h, --help                                                       Help for "stackit config set"
      --iaas-custom-endpoint string                                IaaS API base URL, used in calls to this API
      --identity-provider-custom-client-id string                  Identity Provider client ID, used for user authentication
      --identity-provider-custom-well-known-configuration string   Identity Provider well-known OpenID configuration URL, used for user authentication
      --intake-custom-endpoint string                              Intake API base URL, used in calls to this API
      --kms-custom-endpoint string                                 KMS API base URL, used in calls to this API
      --load-balancer-custom-endpoint string                       Load Balancer API base URL, used in calls to this API
      --logme-custom-endpoint string                               LogMe API base URL, used in calls to this API
      --logs-custom-endpoint string                                Logs API base URL, used in calls to this API
      --mariadb-custom-endpoint string                             MariaDB API base URL, used in calls to this API
      --mongodbflex-custom-endpoint string                         MongoDB Flex API base URL, used in calls to this API
      --object-storage-custom-endpoint string                      Object Storage API base URL, used in calls to this API
      --observability-custom-endpoint string                       Observability API base URL, used in calls to this API
      --opensearch-custom-endpoint string                          OpenSearch API base URL, used in calls to this API
      --postgresflex-custom-endpoint string                        PostgreSQL Flex API base URL, used in calls to this API
      --rabbitmq-custom-endpoint string                            RabbitMQ API base URL, used in calls to this API
      --redis-custom-endpoint string                               Redis API base URL, used in calls to this API
      --resource-manager-custom-endpoint string                    Resource Manager API base URL, used in calls to this API
      --runcommand-custom-endpoint string                          Run Command API base URL, used in calls to this API
      --secrets-manager-custom-endpoint string                     Secrets Manager API base URL, used in calls to this API
      --server-osupdate-custom-endpoint string                     Server Update Management API base URL, used in calls to this API
      --serverbackup-custom-endpoint string                        Server Backup API base URL, used in calls to this API
      --service-account-custom-endpoint string                     Service Account API base URL, used in calls to this API
      --service-enablement-custom-endpoint string                  Service Enablement API base URL, used in calls to this API
      --session-time-limit string                                  Maximum time before authentication is required again. After this time, you will be prompted to login again to execute commands that require authentication. Can't be larger than 24h. Requires authentication after being set to take effect. Examples: 3h, 5h30m40s
      --sfs-custom-endpoint string                                 SFS API base URL, used in calls to this API
      --ske-custom-endpoint string                                 SKE API base URL, used in calls to this API
      --sqlserverflex-custom-endpoint string                       SQLServer Flex API base URL, used in calls to this API
      --token-custom-endpoint string                               Custom token endpoint of the Service Account API, which is used to request access tokens when the service account authentication is activated. Not relevant for user authentication.
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - Provides functionality for CLI configuration options

