## stackit beta server resize

Resizes the server to the given machine type

### Synopsis

Resizes the server to the given machine type.

```
stackit beta server resize [flags]
```

### Examples

```
  Resize a server with ID "xxx" to machine type "yyy"
  $ stackit beta server resize xxx --machine-type yyy
```

### Options

```
  -h, --help                  Help for "stackit beta server resize"
      --machine-type string   Name of the type of the machine for the server. Possible values are documented in https://docs.stackit.cloud/stackit/en/virtual-machine-flavors-75137231.html
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

* [stackit beta server](./stackit_beta_server.md)	 - Provides functionality for servers

