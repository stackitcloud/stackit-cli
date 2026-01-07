## stackit beta sfs snapshot create

Creates a new snapshot of a resource pool

### Synopsis

Creates a new snapshot of a resource pool.

```
stackit beta sfs snapshot create [flags]
```

### Examples

```
  Create a new snapshot with name "snapshot-name" of a resource pool with ID "xxx"
  $ stackit beta sfs snapshot create --name snapshot-name --resource-pool-id xxx

  Create a new snapshot with name "snapshot-name" and comment "snapshot-comment" of a resource pool with ID "xxx"
  $ stackit beta sfs snapshot create --name snapshot-name --resource-pool-id xxx --comment "snapshot-comment"
```

### Options

```
      --comment string            A comment to add more information to the snapshot
  -h, --help                      Help for "stackit beta sfs snapshot create"
      --name string               Snapshot name
      --resource-pool-id string   The resource pool from which the snapshot should be created
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

