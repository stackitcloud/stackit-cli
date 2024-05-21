## stackit argus credentials create

Creates credentials for an Argus instance.

### Synopsis

Creates credentials (username and password) for an Argus instance.
The credentials will be generated and included in the response. You won't be able to retrieve the password later.

```
stackit argus credentials create [flags]
```

### Examples

```
  Create credentials for Argus instance with ID "xxx"
  $ stackit argus credentials create --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit argus credentials create"
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

* [stackit argus credentials](./stackit_argus_credentials.md)	 - Provides functionality for Argus credentials

