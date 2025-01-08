## stackit redis instance list

Lists all Redis instances

### Synopsis

Lists all Redis instances.

```
stackit redis instance list [flags]
```

### Examples

```
  List all Redis instances
  $ stackit redis instance list

  List all Redis instances in JSON format
  $ stackit redis instance list --output-format json

  List up to 10 Redis instances
  $ stackit redis instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit redis instance list"
      --limit int   Maximum number of entries to list
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

* [stackit redis instance](./stackit_redis_instance.md)	 - Provides functionality for Redis instances

