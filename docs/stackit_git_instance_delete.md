## stackit git instance delete

Deletes STACKIT Git instance

### Synopsis

Deletes a STACKIT Git instance by its internal ID.

```
stackit git instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a instance with ID "xxx"
  $ stackit git instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit git instance delete"
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

* [stackit git instance](./stackit_git_instance.md)	 - Provides functionality for STACKIT Git instances

