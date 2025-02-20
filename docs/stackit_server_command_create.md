## stackit server command create

Creates a Server Command

### Synopsis

Creates a Server Command.

```
stackit server command create [flags]
```

### Examples

```
  Create a server command for server with ID "xxx", template name "RunShellScript" and a script from a file (using the @{...} format)
  $ stackit server command create --server-id xxx --template-name=RunShellScript --params script='@{/path/to/script.sh}'

  Create a server command for server with ID "xxx", template name "RunShellScript" and a script provided on the command line
  $ stackit server command create --server-id xxx --template-name=RunShellScript --params script='echo hello'
```

### Options

```
  -h, --help                    Help for "stackit server command create"
  -r, --params stringToString   Params can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
  -s, --server-id string        Server ID
  -n, --template-name string    Template name
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

* [stackit server command](./stackit_server_command.md)	 - Provides functionality for Server Command

