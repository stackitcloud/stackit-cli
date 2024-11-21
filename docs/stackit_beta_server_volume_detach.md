## stackit beta server volume detach

Detaches a volume from a server

### Synopsis

Detaches a volume from a server.

```
stackit beta server volume detach [flags]
```

### Examples

```
  Detaches a volume with ID "xxx" from a server with ID "yyy"
  $ stackit beta server volume detach xxx --server-id yyy
```

### Options

```
  -h, --help               Help for "stackit beta server volume detach"
      --server-id string   Server ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta server volume](./stackit_beta_server_volume.md)	 - Provides functionality for Server volumes

