## stackit redis credentials delete

Deletes credentials of a Redis instance

### Synopsis

Deletes credentials of a Redis instance.

```
stackit redis credentials delete CREDENTIALS_ID [flags]
```

### Examples

```
  Delete credentials with ID "xxx" of Redis instance with ID "yyy"
  $ stackit redis credentials delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit redis credentials delete"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit redis credentials](./stackit_redis_credentials.md)	 - Provides functionality for Redis credentials

