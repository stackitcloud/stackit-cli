## stackit redis instance delete

Deletes a Redis instance

### Synopsis

Deletes a Redis instance.

```
stackit redis instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a Redis instance with ID "xxx"
  $ stackit redis instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit redis instance delete"
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

