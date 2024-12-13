## stackit beta server console

Gets a URL for server remote console

### Synopsis

Gets a URL for server remote console.

```
stackit beta server console [flags]
```

### Examples

```
  Get a URL for the server remote console with server ID "xxx"
  $ stackit beta server console xxx

  Get a URL for the server remote console with server ID "xxx" in JSON format
  $ stackit beta server console xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta server console"
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

