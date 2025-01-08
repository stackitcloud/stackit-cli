## stackit beta server volume update

Updates an attached volume of a server

### Synopsis

Updates an attached volume of a server.

```
stackit beta server volume update VOLUME_ID [flags]
```

### Examples

```
  Update a volume with ID "xxx" of a server with ID "yyy" and enables delete on termination
  $ stackit beta server volume update xxx --server-id yyy --delete-on-termination
```

### Options

```
  -b, --delete-on-termination   Delete the volume during the termination of the server. (default false)
  -h, --help                    Help for "stackit beta server volume update"
      --server-id string        Server ID
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

* [stackit beta server volume](./stackit_beta_server_volume.md)	 - Provides functionality for server volumes

