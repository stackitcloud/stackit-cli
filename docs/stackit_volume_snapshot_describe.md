## stackit volume snapshot describe

Describes a snapshot

### Synopsis

Describes a snapshot by its ID.

```
stackit volume snapshot describe SNAPSHOT_ID [flags]
```

### Examples

```
  Get details of a snapshot with ID "xxx"
  $ stackit volume snapshot describe xxx

  Get details of a snapshot with ID "xxx" in JSON format
  $ stackit volume snapshot describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit volume snapshot describe"
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

