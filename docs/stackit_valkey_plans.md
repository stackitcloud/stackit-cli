## stackit valkey plans

Lists all Valkey service plans

### Synopsis

Lists all Valkey service plans.

```
stackit valkey plans [flags]
```

### Examples

```
  Lists all Valkey service plans
  $ stackit valkey plans

  List all Valkey service plans in JSON format
  $ stackit valkey plans --output-format json

  List up to 10 Valkey service plans
  $ stackit valkey plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit valkey plans"
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

* [stackit valkey](./stackit_valkey.md)	 - Provides functionality for Valkey

