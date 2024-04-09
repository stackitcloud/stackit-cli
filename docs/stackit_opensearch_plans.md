## stackit opensearch plans

Lists all OpenSearch service plans

### Synopsis

Lists all OpenSearch service plans.

```
stackit opensearch plans [flags]
```

### Examples

```
  List all OpenSearch service plans
  $ stackit opensearch plans

  List all OpenSearch service plans in JSON format
  $ stackit opensearch plans --output-format json

  List up to 10 OpenSearch service plans
  $ stackit opensearch plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit opensearch plans"
      --limit int   Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit opensearch](./stackit_opensearch.md)	 - Provides functionality for OpenSearch

