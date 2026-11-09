## stackit container-registry artifactory list

Lists all Container Registry artifactories

### Synopsis

Lists all Container Registry artifactories.

```
stackit container-registry artifactory list [flags]
```

### Examples

```
  List all artifactories
  $ stackit container-registry artifactory list

  List all artifactories in JSON format
  $ stackit container-registry artifactory list --output-format json

  List up to 10 artifactories
  $ stackit container-registry artifactory list --limit 10
```

### Options

```
  -h, --help        Help for "stackit container-registry artifactory list"
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

* [stackit container-registry artifactory](./stackit_container-registry_artifactory.md)	 - Provides functionality for Container Registry artifactories

