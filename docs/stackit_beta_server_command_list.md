## stackit beta server command list

Lists all server commands

### Synopsis

Lists all server commands.

```
stackit beta server command list [flags]
```

### Examples

```
  List all commands for a server with ID "xxx"
  $ stackit beta server command list --server-id xxx

  List all commands for a server with ID "xxx" in JSON format
  $ stackit beta server command list --server-id xxx --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server command list"
      --limit int          Maximum number of entries to list
  -s, --server-id string   Server ID
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

* [stackit beta server command](./stackit_beta_server_command.md)	 - Provides functionality for Server Command

