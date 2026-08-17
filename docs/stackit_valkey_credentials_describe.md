## stackit valkey credentials describe

Shows details of credentials of a Valkey instance

### Synopsis

Shows details of credentials of a Valkey instance. The password will be shown in plain text in the output.

```
stackit valkey credentials describe CREDENTIALS_ID [flags]
```

### Examples

```
  Get details of credentials with ID "xxx" from a Valkey instance with ID "yyy"
  $ stackit beta valkey credentials describe xxx --instance-id yyy

  Get details of credentials with ID "xxx" from a Valkey instance with ID "yyy" in JSON format
  $ stackit beta valkey credentials describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help                 Help for "stackit valkey credentials describe"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit valkey credentials](./stackit_valkey_credentials.md)	 - Provides functionality for Valkey credentials

