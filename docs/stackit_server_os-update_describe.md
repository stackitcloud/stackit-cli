## stackit server os-update describe

Shows details of a Server os-update

### Synopsis

Shows details of a Server os-update.

```
stackit server os-update describe UPDATE_ID [flags]
```

### Examples

```
  Get details of a Server os-update with id "my-os-update-id"
  $ stackit server os-update describe my-os-update-id

  Get details of a Server os-update with id "my-os-update-id" in JSON format
  $ stackit server os-update describe my-os-update-id --output-format json
```

### Options

```
  -h, --help               Help for "stackit server os-update describe"
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

* [stackit server os-update](./stackit_server_os-update.md)	 - Provides functionality for managed server updates

