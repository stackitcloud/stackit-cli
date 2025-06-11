## stackit volume snapshot delete

Deletes a snapshot

### Synopsis

Deletes a snapshot by its ID.

```
stackit volume snapshot delete SNAPSHOT_ID [flags]
```

### Examples

```
  Delete a snapshot with ID "xxx"
  $ stackit volume snapshot delete xxx
```

### Options

```
  -h, --help   Help for "stackit volume snapshot delete"
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

