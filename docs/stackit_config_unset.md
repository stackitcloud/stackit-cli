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
      --argus-custom-endpoint              Argus API base URL. If unset, uses the default base URL
      --async                              Configuration option to run commands asynchronously
      --authorization-custom-endpoint      Authorization API base URL. If unset, uses the default base URL
      --dns-custom-endpoint                DNS API base URL. If unset, uses the default base URL
  -h, --help                               Help for "stackit config unset"
      --load-balancer-custom-endpoint      Load Balancer API base URL. If unset, uses the default base URL
      --logme-custom-endpoint              LogMe API base URL. If unset, uses the default base URL
      --mariadb-custom-endpoint            MariaDB API base URL. If unset, uses the default base URL
      --mongodbflex-custom-endpoint        MongoDB Flex API base URL. If unset, uses the default base URL
      --object-storage-custom-endpoint     Object Storage API base URL. If unset, uses the default base URL
      --opensearch-custom-endpoint         OpenSearch API base URL. If unset, uses the default base URL
      --output-format                      Output format
      --postgresflex-custom-endpoint       PostgreSQL Flex API base URL. If unset, uses the default base URL
      --project-id                         Project ID
      --rabbitmq-custom-endpoint           RabbitMQ API base URL. If unset, uses the default base URL
      --redis-custom-endpoint              Redis API base URL. If unset, uses the default base URL
      --resource-manager-custom-endpoint   Resource Manager API base URL. If unset, uses the default base URL
      --secrets-manager-custom-endpoint    Secrets Manager API base URL. If unset, uses the default base URL
      --service-account-custom-endpoint    SKE API base URL. If unset, uses the default base URL
      --session-time-limit                 Maximum time before authentication is required again. If unset, defaults to 2h
      --ske-custom-endpoint                SKE API base URL. If unset, uses the default base URL
      --verbosity                          Verbosity of the CLI
```

### Options inherited from parent commands

```
  -y, --assume-yes   If set, skips all confirmation prompts
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - Provides functionality for CLI configuration options

