## stackit beta server reboot

Reboots a server

### Synopsis

Reboots a server.

```
stackit beta server reboot SERVER_ID [flags]
```

### Examples

```
  Perform a soft reboot of a server with ID "xxx"
  $ stackit beta server reboot xxx

  Perform a hard reboot of a server with ID "xxx"
  $ stackit beta server reboot xxx --hard
```

### Options

```
  -b, --hard   Performs a hard reboot. (default false)
  -h, --help   Help for "stackit beta server reboot"
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

