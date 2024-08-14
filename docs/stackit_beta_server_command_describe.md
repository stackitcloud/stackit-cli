## stackit beta server command describe

Shows details of a Server Command

### Synopsis

Shows details of a Server Command.

```
stackit beta server command describe COMMAND_ID [flags]
```

### Examples

```
  Get details of a Server Command with ID "xxx" for server with ID "yyy"
  $ stackit beta server command describe xxx --server-id=yyy

  Get details of a Server Command with ID "xxx" for server with ID "yyy" in JSON format
  $ stackit beta server command describe xxx --server-id=yyy --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server command describe"
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

