## stackit network list

Lists all networks of a project

### Synopsis

Lists all network of a project.

```
stackit network list [flags]
```

### Examples

```
  Lists all networks
  $ stackit network list

  Lists all networks in JSON format
  $ stackit network list --output-format json

  Lists up to 10 networks
  $ stackit network list --limit 10

  Lists all networks which contains the label xxx
  $ stackit network list --label-selector xxx
```

### Options

```
  -h, --help                    Help for "stackit network list"
      --label-selector string   Filter by label
      --limit int               Maximum number of entries to list
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

* [stackit network](./stackit_network.md)	 - Provides functionality for networks

