## stackit beta server unrescue

Unrescues an existing server

### Synopsis

Unrescues an existing server.

```
stackit beta server unrescue SERVER_ID [flags]
```

### Examples

```
  Unrescue an existing server with ID "xxx"
  $ stackit beta server unrescue xxx
```

### Options

```
  -h, --help   Help for "stackit beta server unrescue"
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

* [stackit beta server](./stackit_beta_server.md)	 - Provides functionality for servers

