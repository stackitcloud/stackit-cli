## stackit beta affinity-group delete

Deletes an affinity group

### Synopsis

Deletes an affinity group.

```
stackit beta affinity-group delete AFFINITY_GROUP [flags]
```

### Examples

```
  Delete an affinity group with ID "xxx"
  $ stackit beta affinity-group delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta affinity-group delete"
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

* [stackit beta affinity-group](./stackit_beta_affinity-group.md)	 - Manage server affinity groups

