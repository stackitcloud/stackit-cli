## stackit redis plans

Lists all Redis service plans

### Synopsis

Lists all Redis service plans.

```
stackit redis plans [flags]
```

### Examples

```
  List all Redis service plans
  $ stackit redis plans

  List all Redis service plans in JSON format
  $ stackit redis plans --output-format json

  List up to 10 Redis service plans
  $ stackit redis plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit redis plans"
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

* [stackit redis](./stackit_redis.md)	 - Provides functionality for Redis

