## stackit config set

Sets CLI configuration options

### Synopsis

Sets CLI configuration options.

```
stackit config set [flags]
```

### Examples

```
  Set a project ID in your active configuration. This project ID will be used by every command (as long as it's not overridden by the "STACKIT_PROJECT_ID" environment variable)
  $ stackit config set --project-id xxx

  Set the session time limit to 1 hour
  $ stackit config set --session-time-limit 1h

  Set the DNS custom endpoint
  $ stackit config set --dns-custom-endpoint https://dns.stackit.cloud
```

### Options

```
      --authorization-custom-endpoint string      Authorization custom endpoint
      --dns-custom-endpoint string                DNS custom endpoint. Will be used as the base URL on all calls to this API
  -h, --help                                      Help for "stackit config set"
      --logme-custom-endpoint string              LogMe custom endpoint. Will be used as the base URL on all calls to this API
      --mariadb-custom-endpoint string            MariaDB custom endpoint. Will be used as the base URL on all calls to this API
      --mongodbflex-custom-endpoint string        MongoDB Flex custom endpoint. Will be used as the base URL on all calls to this API
      --opensearch-custom-endpoint string         OpenSearch custom endpoint. Will be used as the base URL on all calls to this API
      --postgresflex-custom-endpoint string       PostgreSQL Flex custom endpoint. Will be used as the base URL on all calls to this API
      --rabbitmq-custom-endpoint string           RabbitMQ custom endpoint. Will be used as the base URL on all calls to this API
      --redis-custom-endpoint string              Redis custom endpoint. Will be used as the base URL on all calls to this API
      --resource-manager-custom-endpoint string   Resource manager custom endpoint. Will be used as the base URL on all calls to this API
      --service-account-custom-endpoint string    Service Account custom endpoint. Will be used as the base URL on all calls to this API
      --session-time-limit string                 Maximum time before authentication is required again. After this time, you will be prompted to login again to execute commands that require authentication. Can't be larger than 24h. Requires authentication after being set to take effect. Examples: 3h, 5h30m40s (BETA: currently values greater than 2h have no effect)
      --ske-custom-endpoint string                SKE custom endpoint. Will be used as the base URL on all calls to this API
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - Provides functionality for CLI configuration options

