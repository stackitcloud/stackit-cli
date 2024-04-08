## stackit logme credentials create

Creates credentials for a LogMe instance

### Synopsis

Creates credentials (username and password) for a LogMe instance.

```
stackit logme credentials create [flags]
```

### Examples

```
  Create credentials for a LogMe instance
  $ stackit logme credentials create --instance-id xxx

  Create credentials for a LogMe instance and hide the password in the output
  $ stackit logme credentials create --instance-id xxx --hide-password
```

### Options

```
  -h, --help                 Help for "stackit logme credentials create"
      --hide-password        Hide password in output
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit logme credentials](./stackit_logme_credentials.md)	 - Provides functionality for LogMe credentials

