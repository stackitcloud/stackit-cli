## stackit beta server log

Gets server console log

### Synopsis

Gets server console log.

```
stackit beta server log SERVER_ID [flags]
```

### Examples

```
  Get server console log for the server with ID "xxx"
  $ stackit beta server log xxx

  Get server console log for the server with ID "xxx" and limit output lines to 1000
  $ stackit beta server log xxx --length 1000

  Get server console log for the server with ID "xxx" in JSON format
  $ stackit beta server log xxx --output-format json
```

### Options

```
  -h, --help         Help for "stackit beta server log"
      --length int   Maximum number of lines to list (default 2000)
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

