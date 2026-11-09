## stackit config unset

Unsets CLI configuration options

### Synopsis

Unsets CLI configuration options, undoing past usages of the `stackit config set` command.

```
stackit config unset [flags]
```

### Examples

```
  Unset the project ID stored in your configuration
  $ stackit config unset --project-id

  Unset the session time limit stored in your configuration
  $ stackit config unset --session-time-limit

  Unset the DNS custom endpoint stored in your configuration
  $ stackit config unset --dns-custom-endpoint
```

### Options

```
      --alb-waf-custom-endpoint[=true|false]                             ALB WAF API base URL. If unset, uses the default base URL
      --allowed-url-domain[=true|false]                                  Domain name, used for the verification of the URLs that are given in the IDP endpoint and curl commands. If unset, defaults to stackit.cloud
      --assume-yes[=true|false]                                          If set, skips all confirmation prompts
      --async[=true|false]                                               Configuration option to run commands asynchronously
      --authorization-custom-endpoint[=true|false]                       Authorization API base URL. If unset, uses the default base URL
      --automation-custom-endpoint[=true|false]                          Automation API base URL. If unset, uses the default base URL
      --cdn-custom-endpoint[=true|false]                                 Custom CDN endpoint URL. If unset, uses the default base URL
      --dns-custom-endpoint[=true|false]                                 DNS API base URL. If unset, uses the default base URL
      --edge-custom-endpoint[=true|false]                                Edge API base URL. If unset, uses the default base URL
  -h, --help[=true|false]                                                Help for "stackit config unset"
      --iaas-alpha-custom-endpoint[=true|false]                          IaaS Alpha API base URL. If unset, uses the default base URL
      --iaas-custom-endpoint[=true|false]                                IaaS API base URL. If unset, uses the default base URL
      --identity-provider-custom-client-id[=true|false]                  Identity Provider client ID, used for user authentication
      --identity-provider-custom-well-known-configuration[=true|false]   Identity Provider well-known OpenID configuration URL. If unset, uses the default identity provider
      --intake-custom-endpoint[=true|false]                              Intake API base URL. If unset, uses the default base URL
      --kms-custom-endpoint[=true|false]                                 KMS API base URL. If unset, uses the default base URL
      --load-balancer-custom-endpoint[=true|false]                       Load Balancer API base URL. If unset, uses the default base URL
      --logme-custom-endpoint[=true|false]                               LogMe API base URL. If unset, uses the default base URL
      --logs-custom-endpoint[=true|false]                                Logs API base URL. If unset, uses the default base URL
      --mariadb-custom-endpoint[=true|false]                             MariaDB API base URL. If unset, uses the default base URL
      --mongodbflex-custom-endpoint[=true|false]                         MongoDB Flex API base URL. If unset, uses the default base URL
      --object-storage-custom-endpoint[=true|false]                      Object Storage API base URL. If unset, uses the default base URL
      --observability-custom-endpoint[=true|false]                       Observability API base URL. If unset, uses the default base URL
      --opensearch-custom-endpoint[=true|false]                          OpenSearch API base URL. If unset, uses the default base URL
      --output-format[=true|false]                                       Output format
      --postgresflex-custom-endpoint[=true|false]                        PostgreSQL Flex API base URL. If unset, uses the default base URL
      --project-id[=true|false]                                          Project ID
      --rabbitmq-custom-endpoint[=true|false]                            RabbitMQ API base URL. If unset, uses the default base URL
      --redis-custom-endpoint[=true|false]                               Redis API base URL. If unset, uses the default base URL
      --region[=true|false]                                              Region
      --resource-manager-custom-endpoint[=true|false]                    Resource Manager API base URL. If unset, uses the default base URL
      --runcommand-custom-endpoint[=true|false]                          Server Command base URL. If unset, uses the default base URL
      --scf-custom-endpoint[=true|false]                                 STACKIT Cloud Foundry API base URL. If unset, uses the default base URL
      --secrets-manager-custom-endpoint[=true|false]                     Secrets Manager API base URL. If unset, uses the default base URL
      --server-osupdate-custom-endpoint[=true|false]                     Server Update Management base URL. If unset, uses the default base URL
      --serverbackup-custom-endpoint[=true|false]                        Server Backup base URL. If unset, uses the default base URL
      --service-account-custom-endpoint[=true|false]                     Service Account API base URL. If unset, uses the default base URL
      --service-enablement-custom-endpoint[=true|false]                  Service Enablement API base URL. If unset, uses the default base URL
      --session-time-limit[=true|false]                                  Maximum time before authentication is required again. If unset, defaults to 12h
      --sfs-custom-endpoint[=true|false]                                 SFS API base URL. If unset, uses the default base URL
      --ske-custom-endpoint[=true|false]                                 SKE API base URL. If unset, uses the default base URL
      --sqlserverflex-custom-endpoint[=true|false]                       SQLServer Flex API base URL. If unset, uses the default base URL
      --token-custom-endpoint[=true|false]                               Custom token endpoint of the Service Account API, which is used to request access tokens when the service account authentication is activated. Not relevant for user authentication.
      --verbosity[=true|false]                                           Verbosity of the CLI
      --vpn-custom-endpoint[=true|false]                                 VPN API base URL. If unset, uses the default base URL
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - Provides functionality for CLI configuration options

