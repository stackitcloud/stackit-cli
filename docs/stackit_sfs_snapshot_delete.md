## stackit sfs snapshot delete

Deletes a snapshot

### Synopsis

Deletes a snapshot.

```
stackit sfs snapshot delete SNAPSHOT_NAME [flags]
```

### Examples

```
  Delete a snapshot with "SNAPSHOT_NAME" from resource pool with ID "yyy"
  $ stackit sfs snapshot delete SNAPSHOT_NAME --resource-pool-id yyy
```

### Options

```
  -h, --help                      Help for "stackit sfs snapshot delete"
      --resource-pool-id string   The resource pool from which the snapshot should be created
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

* [stackit sfs snapshot](./stackit_sfs_snapshot.md)	 - Provides functionality for SFS snapshots

