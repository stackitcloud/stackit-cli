## stackit beta valkey instance delete

Deletes a Valkey instance

### Synopsis

Deletes a Valkey instance.

```
stackit beta valkey instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a Valkey instance with ID "xxx"
  $ stackit beta valkey instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta valkey instance delete"
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

* [stackit beta valkey instance](./stackit_beta_valkey_instance.md)	 - Provides functionality for Valkey instances

