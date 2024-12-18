## stackit beta server start

Starts an existing server or allocates the server if deallocated

### Synopsis

Starts an existing server or allocates the server if deallocated.

```
stackit beta server start [flags]
```

### Examples

```
  Start an existing server with ID "xxx"
  $ stackit beta server start xxx
```

### Options

```
  -h, --help   Help for "stackit beta server start"
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

* [stackit beta server](./stackit_beta_server.md)	 - Provides functionality for servers
