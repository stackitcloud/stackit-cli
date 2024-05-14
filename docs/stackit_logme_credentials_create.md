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

  Create credentials for a LogMe instance and show the password in the output
  $ stackit logme credentials create --instance-id xxx --show-password
```

### Options

```
  -h, --help                 Help for "stackit logme credentials create"
      --instance-id string   Instance ID
  -s, --show-password        Show password in output
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

* [stackit logme credentials](./stackit_logme_credentials.md)	 - Provides functionality for LogMe credentials

