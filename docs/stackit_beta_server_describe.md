## stackit beta server describe

Shows details of a server

### Synopsis

Shows details of a server.

```
stackit beta server describe SERVER_ID [flags]
```

### Examples

```
  Show details of a server with ID "xxx"
  $ stackit beta server describe xxx

  Show details of a server with ID "xxx" in JSON format
  $ stackit beta server describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta server describe"
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

