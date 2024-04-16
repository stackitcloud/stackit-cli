## stackit argus credentials create

Creates credentials for an Argus instance.

### Synopsis

Creates credentials for an Argus instance.

```
stackit argus credentials create [flags]
```

### Examples

```
  Create credentials for Argus instance with ID "xxx"
  $ stackit argus credentials create --instance-id xxx

  Create credentials for Argus instance with ID "xxx" and hide the password in the output
  $ stackit argus credentials create --instance-id xxx --hide-password
```

### Options

```
  -h, --help                 Help for "stackit argus credentials create"
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

* [stackit argus credentials](./stackit_argus_credentials.md)	 - Provides functionality for Argus credentials

