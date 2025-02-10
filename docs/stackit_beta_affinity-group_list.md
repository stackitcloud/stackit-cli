## stackit beta affinity-group list

Lists affinity groups

### Synopsis

Lists affinity groups.

```
stackit beta affinity-group list [flags]
```

### Examples

```
  Lists all affinity groups
  $ stackit beta affinity-group list

  Lists up to 10 affinity groups
  $ stackit beta affinity-group list --limit=10
```

### Options

```
  -h, --help        Help for "stackit beta affinity-group list"
      --limit int   Limit the output to the first n elements
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

