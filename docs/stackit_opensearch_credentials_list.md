## stackit opensearch credentials list

Lists all credentials' IDs for an OpenSearch instance

### Synopsis

Lists all credentials' IDs for an OpenSearch instance.

```
stackit opensearch credentials list [flags]
```

### Examples

```
  List all credentials' IDs for an OpenSearch instance
  $ stackit opensearch credentials list --instance-id xxx

  List all credentials' IDs for an OpenSearch instance in JSON format
  $ stackit opensearch credentials list --instance-id xxx --output-format json

  List up to 10 credentials' IDs for an OpenSearch instance
  $ stackit opensearch credentials list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit opensearch credentials list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit opensearch credentials](./stackit_opensearch_credentials.md)	 - Provides functionality for OpenSearch credentials

