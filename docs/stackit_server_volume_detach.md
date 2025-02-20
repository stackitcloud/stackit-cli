## stackit server volume detach

Detaches a volume from a server

### Synopsis

Detaches a volume from a server.

```
stackit server volume detach VOLUME_ID [flags]
```

### Examples

```
  Detaches a volume with ID "xxx" from a server with ID "yyy"
  $ stackit server volume detach xxx --server-id yyy
```

### Options

```
  -h, --help               Help for "stackit server volume detach"
      --server-id string   Server ID
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

* [stackit server volume](./stackit_server_volume.md)	 - Provides functionality for server volumes

