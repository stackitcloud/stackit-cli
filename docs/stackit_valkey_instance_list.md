## stackit valkey instance list

Lists all Valkey instances

### Synopsis

Lists all Valkey instances.

```
stackit valkey instance list [flags]
```

### Examples

```
  List all Valkey instances
  $ stackit valkey instance list

  List all Valkey instances in JSON format
  $ stackit valkey instance list --output-format json

  List up to 10 Valkey instances
  $ stackit valkey instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit valkey instance list"
      --limit int   Maximum number of entries to list
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

* [stackit valkey instance](./stackit_valkey_instance.md)	 - Provides functionality for Valkey instances

