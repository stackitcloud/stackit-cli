## stackit beta sfs snapshot update

Updates a new snapshot of a resource pool

### Synopsis

Updates a new snapshot of a resource pool.

```
stackit beta sfs snapshot update SNAPSHOT_NAME [flags]
```

### Examples

```
  Updates the name of a snapshot with name "snapshot-name" of a resource pool with ID "xxx"
  $ stackit beta sfs snapshot update snapshot-name --resource-pool-id xxx --name new-snapshot-name

  Updates the comment of a snapshot with name "snapshot-name" of a resource pool with ID "xxx"
  $ stackit beta sfs snapshot update snapshot-name --resource-pool-id xxx --comment "snapshot-comment"
```

### Options

```
      --comment string            A comment to add more information to the snapshot
  -h, --help                      Help for "stackit beta sfs snapshot update"
      --name string               Snapshot name
      --resource-pool-id string   The resource pool from which the snapshot should be updated
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

* [stackit beta sfs snapshot](./stackit_beta_sfs_snapshot.md)	 - Provides functionality for SFS snapshots

