## stackit sqlserverflex instance create

Creates credentials for an Argus instance.

### Synopsis

Creates credentials (username and password) for an Argus instance.
The credentials will be generated and included in the response. You won't be able to retrieve the password later.

```
stackit sqlserverflex instance create [flags]
```

### Examples

```
  Create credentials for Argus instance with ID "xxx"
  $ stackit argus credentials create --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit sqlserverflex instance create"
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

* [stackit sqlserverflex instance](./stackit_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

