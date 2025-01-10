## stackit beta server service-account list

List all attached service accounts for a server

### Synopsis

List all attached service accounts for a server

```
stackit beta server service-account list [flags]
```

### Examples

```
  List all attached service accounts for a server with ID "xxx"
  $ stackit beta server service-account list --server-id xxx

  List up to 10 attached service accounts for a server with ID "xxx"
  $ stackit beta server service-account list --server-id xxx --limit 10

  List all attached service accounts for a server with ID "xxx" in JSON format
  $ stackit beta server service-account list --server-id xxx --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server service-account list"
      --limit int          Maximum number of entries to list
  -s, --server-id string   Server ID
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

* [stackit beta server service-account](./stackit_beta_server_service-account.md)	 - Allows attaching/detaching service accounts to servers

