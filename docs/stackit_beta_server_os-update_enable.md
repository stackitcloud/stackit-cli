## stackit beta server os-update enable

Enables Server os-update service

### Synopsis

Enables Server os-update service.

```
stackit beta server os-update enable [flags]
```

### Examples

```
  Enable os-update functionality for your server
  $ stackit beta server os-update enable --server-id=zzz
```

### Options

```
  -h, --help               Help for "stackit beta server os-update enable"
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

* [stackit beta server os-update](./stackit_beta_server_os-update.md)	 - Provides functionality for managed server updates

