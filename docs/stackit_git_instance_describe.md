## stackit git instance describe

Describes STACKIT Git instance

### Synopsis

Describes a STACKIT Git instance by its internal ID.

```
stackit git instance describe INSTANCE_ID [flags]
```

### Examples

```
  Describe instance "xxx"
  $ stackit git describe xxx
```

### Options

```
  -h, --help   Help for "stackit git instance describe"
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

