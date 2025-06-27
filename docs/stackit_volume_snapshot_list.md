## stackit volume snapshot list

Lists all snapshots

### Synopsis

Lists all snapshots in a project.

```
stackit volume snapshot list [flags]
```

### Examples

```
  List all snapshots
  $ stackit volume snapshot list

  List snapshots with a limit of 10
  $ stackit volume snapshot list --limit 10

  List snapshots filtered by label
  $ stackit volume snapshot list --label-selector key1=value1
```

### Options

```
  -h, --help                    Help for "stackit volume snapshot list"
      --label-selector string   Filter snapshots by labels
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

* [stackit volume snapshot](./stackit_volume_snapshot.md)	 - Provides functionality for snapshots

