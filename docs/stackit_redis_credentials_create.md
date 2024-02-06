## stackit redis credentials create

Creates credentials for a Redis instance

### Synopsis

Creates credentials (username and password) for a Redis instance.

```
stackit redis credentials create [flags]
```

### Examples

```
  Create credentials for a Redis instance
  $ stackit redis credentials create --instance-id xxx

  Create credentials for a Redis instance and hide the password in the output
  $ stackit redis credentials create --instance-id xxx --hide-password
```

### Options

```
  -h, --help                 Help for "stackit redis credentials create"
      --hide-password        Hide password in output
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit redis credentials](./stackit_redis_credentials.md)	 - Provides functionality for Redis credentials

