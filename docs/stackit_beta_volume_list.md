## stackit beta volume list

Lists all volumes of a project

### Synopsis

Lists all volumes of a project.

```
stackit beta volume list [flags]
```

### Examples

```
  Lists all volumes
  $ stackit beta volume list

  Lists all volumes which contains the label xxx
  $ stackit beta volume list --label-selector xxx

  Lists all volumes in JSON format
  $ stackit beta volume list --output-format json

  Lists up to 10 volumes
  $ stackit beta volume list --limit 10
```

### Options

```
  -h, --help                    Help for "stackit beta volume list"
      --label-selector string   Filter by label
      --limit int               Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for Volume
