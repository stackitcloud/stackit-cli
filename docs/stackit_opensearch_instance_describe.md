## stackit opensearch instance describe

Get details of an OpenSearch instance

### Synopsis

Get details of an OpenSearch instance.

```
stackit opensearch instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of an OpenSearch instance with ID "xxx"
  $ stackit opensearch instance describe xxx

  Get details of an OpenSearch instance with ID "xxx" in a table format
  $ stackit opensearch instance describe xxx --output-format pretty
```

### Options

```
  -h, --help   Help for "stackit opensearch instance describe"
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

