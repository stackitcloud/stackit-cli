## stackit opensearch instance list

List all OpenSearch instances

### Synopsis

List all OpenSearch instances

```
stackit opensearch instance list [flags]
```

### Examples

```
  List all OpenSearch instances
  $ stackit opensearch instance list

  List all OpenSearch instances in JSON format
  $ stackit opensearch instance list --output-format json

  List up to 10 OpenSearch instances
  $ stackit opensearch instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit opensearch instance list"
      --limit int   Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit opensearch instance](./stackit_opensearch_instance.md)	 - Provides functionality for OpenSearch instances

