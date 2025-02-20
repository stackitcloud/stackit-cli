## stackit server rescue

Rescues an existing server

### Synopsis

Rescues an existing server.

```
stackit server rescue SERVER_ID [flags]
```

### Examples

```
  Rescue an existing server with ID "xxx" using image with ID "yyy" as boot volume
  $ stackit server rescue xxx --image-id yyy
```

### Options

```
  -h, --help              Help for "stackit server rescue"
      --image-id string   The image ID to be used for a temporary boot volume.
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

