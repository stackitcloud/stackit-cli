## stackit beta server update

Updates a server

### Synopsis

Updates a server.

```
stackit beta server update [flags]
```

### Examples

```
  Update server with ID "xxx" with new name "server-1-new"
  $ stackit beta server update xxx --name server-1-new

  Update server with ID "xxx" with new name "server-1-new" and label(s)
  $ stackit beta server update xxx --name server-1-new --labels key=value,foo=bar
```

### Options

```
  -h, --help                    Help for "stackit beta server update"
      --labels stringToString   Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...' (default [])
  -n, --name string             Server name
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

* [stackit beta server](./stackit_beta_server.md)	 - Provides functionality for Server
