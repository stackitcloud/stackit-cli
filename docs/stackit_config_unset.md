## stackit config unset

Unsets CLI configuration options

### Synopsis

Unsets CLI configuration options.

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
      --async                              Configuration option to run commands asynchronously
      --dns-custom-endpoint                DNS custom endpoint
  -h, --help                               Help for "stackit config unset"
      --membership-custom-endpoint         Membership custom endpoint
      --mongodbflex-custom-endpoint        MongoDB Flex custom endpoint
      --opensearch-custom-endpoint         OpenSearch custom endpoint
      --output-format                      Output format
      --project-id                         Project ID
      --resource-manager-custom-endpoint   Resource Manager custom endpoint
      --service-account-custom-endpoint    SKE custom endpoint
      --ske-custom-endpoint                SKE custom endpoint
```

### Options inherited from parent commands

```
  -y, --assume-yes   If set, skips all confirmation prompts
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - CLI configuration options

