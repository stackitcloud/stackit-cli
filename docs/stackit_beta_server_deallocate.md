## stackit beta server deallocate

Deallocates an existing server

### Synopsis

Deallocates an existing server.

```
stackit beta server deallocate SERVER_ID [flags]
```

### Examples

```
  Deallocate an existing server with ID "xxx"
  $ stackit beta server deallocate xxx
```

### Options

```
  -h, --help   Help for "stackit beta server deallocate"
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

