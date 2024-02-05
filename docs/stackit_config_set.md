## stackit config set

Sets CLI configuration options

### Synopsis

Sets CLI configuration options.
All of the configuration options can be set using an environment variable, which takes precedence over what is configured.
The environment variable is the name of the flag, with underscores ("_") instead of dashes ("-") and the "STACKIT" prefix.
Example: to set the project ID you can set the environment variable STACKIT_PROJECT_ID.

```
stackit config set [flags]
```

### Examples

```
  Set a project ID in your active configuration. This project ID will be used by every command, as long as it's not overridden by the "STACKIT_PROJECT_ID" environment variable or the command flag
  $ stackit config set --project-id xxx

  Set the session time limit to 1 hour. After this time you will be prompted to login again to be able to execute commands that need authentication
  $ stackit config set --session-time-limit 1h

  Set the DNS custom endpoint. This endpoint will be used on all calls to the DNS API, unless overridden by the "STACKIT_DNS_CUSTOM_ENDPOINT" environment variable
  $ stackit config set --dns-custom-endpoint https://dns.stackit.cloud
```

### Options

```
      --dns-custom-endpoint string                DNS custom endpoint
  -h, --help                                      Help for "stackit config set"
      --membership-custom-endpoint string         Membership custom endpoint
      --mongodbflex-custom-endpoint string        MongoDB Flex custom endpoint
      --opensearch-custom-endpoint string         OpenSearch custom endpoint
      --resource-manager-custom-endpoint string   Resource manager custom endpoint
      --service-account-custom-endpoint string    Service Account custom endpoint
      --session-time-limit string                 Maximum time before authentication is required again. Can't be larger than 24h. Examples: 3h, 5h30m40s (BETA: currently values greater than 2h have no effect)
      --ske-custom-endpoint string                SKE custom endpoint
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - CLI configuration options

