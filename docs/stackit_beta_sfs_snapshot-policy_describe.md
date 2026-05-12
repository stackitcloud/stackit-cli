## stackit beta sfs snapshot-policy describe

Shows details of a snapshot policy

### Synopsis

Shows details of a snapshot policy.

```
stackit beta sfs snapshot-policy describe SNAPSHOT_POLICY_ID [flags]
```

### Examples

```
  Describe a snapshot policy with ID "xxx"
  $ stackit beta sfs snapshot-policy describe xxx
```

### Options

```
  -h, --help   Help for "stackit beta sfs snapshot-policy describe"
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

* [stackit beta sfs snapshot-policy](./stackit_beta_sfs_snapshot-policy.md)	 - Provides functionality for SFS snapshot policies

