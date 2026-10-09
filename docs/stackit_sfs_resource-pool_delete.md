## stackit sfs resource-pool delete

Deletes a SFS resource pool

### Synopsis

Deletes a SFS resource pool.

```
stackit sfs resource-pool delete [flags]
```

### Examples

```
  Delete the SFS resource pool with ID "xxx"
  $ stackit sfs resource-pool delete xxx
```

### Options

```
  -h, --help   Help for "stackit sfs resource-pool delete"
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

* [stackit sfs resource-pool](./stackit_sfs_resource-pool.md)	 - Provides functionality for SFS resource pools

