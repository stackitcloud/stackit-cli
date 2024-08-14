## stackit beta server command template describe

Shows details of a Server Command Template

### Synopsis

Shows details of a Server Command Template.

```
stackit beta server command template describe COMMAND_TEMPLATE_NAME [flags]
```

### Examples

```
  Get details of a Server Command Template with name "RunShellScript" for server with ID "xxx"
  $ stackit beta server command template describe RunShellScript --server-id=xxx

  Get details of a Server Command Template with name "RunShellScript" for server with ID "xxx" in JSON format
  $ stackit beta server command template describe RunShellScript --server-id=xxx --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server command template describe"
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

* [stackit beta server command template](./stackit_beta_server_command_template.md)	 - Provides functionality for Server Command Template

