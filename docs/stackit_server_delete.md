## stackit server delete

Deletes a server

### Synopsis

Deletes a server.
If the server is still in use, the deletion will fail


```
stackit server delete SERVER_ID [flags]
```

### Examples

```
  Delete server with ID "xxx"
  $ stackit server delete xxx
```

### Options

```
  -h, --help   Help for "stackit server delete"
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

* [stackit server](./stackit_server.md)	 - Provides functionality for servers

