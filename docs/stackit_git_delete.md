## stackit git delete

Deletes STACKIT Git instance

### Synopsis

Deletes an STACKIT Git instance by its internal ID.

```
stackit git delete INSTANCE_ID [flags]
```

### Examples

```
  Delete an instance with ID "xxx"
  $ stackit git delete xxx
```

### Options

```
  -h, --help   Help for "stackit git delete"
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

* [stackit git](./stackit_git.md)	 - Provides functionality for STACKIT Git

