## stackit valkey credentials create

Creates credentials for a Valkey instance

### Synopsis

Creates credentials (username and password) for a Valkey instance.

```
stackit valkey credentials create [flags]
```

### Examples

```
  Create credentials for a Valkey instance with ID "xxx"
  $ stackit valkey credentials create --instance-id xxx

  Create credentials for a Valkey instance and show the password in the output
  $ stackit valkey credentials create --instance-id xxx --show-password
```

### Options

```
  -h, --help                 Help for "stackit valkey credentials create"
      --instance-id string   Instance ID
  -s, --show-password        Show password in output
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

