## stackit volume snapshot create

Creates a snapshot from a volume

### Synopsis

Creates a snapshot from a volume.

```
stackit volume snapshot create [flags]
```

### Examples

```
  Create a snapshot from a volume with ID "xxx"
  $ stackit volume snapshot create --volume-id xxx

  Create a snapshot from a volume with ID "xxx" and name "my-snapshot"
  $ stackit volume snapshot create --volume-id xxx --name my-snapshot

  Create a snapshot from a volume with ID "xxx" and labels
  $ stackit volume snapshot create --volume-id xxx --labels key1=value1,key2=value2
```

### Options

```
  -h, --help                    Help for "stackit volume snapshot create"
      --labels stringToString   Key-value string pairs as labels (default [])
      --name string             Name of the snapshot
      --volume-id string        Volume ID
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

