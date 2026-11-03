## stackit registry replication-adapter list

Lists all Container Registry replication adapters

### Synopsis

Lists all available adapters for replication to/from other registries.

```
stackit registry replication-adapter list [flags]
```

### Examples

```
  List all replication adapters
  $ stackit registry replication-adapter list

  List all replication adapters in JSON format
  $ stackit registry replication-adapter list --output-format json

  List up to 10 replication adapters
  $ stackit registry replication-adapter list --limit 10
```

### Options

```
  -h, --help        Help for "stackit registry replication-adapter list"
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

* [stackit registry replication-adapter](./stackit_registry_replication-adapter.md)	 - Provides functionality for Container Registry replication adapters

