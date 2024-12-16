## stackit beta server public-ip detach

Detaches a public IP from a server

### Synopsis

Detaches a public IP from a server.

```
stackit beta server public-ip detach PUBLIC_IP_ID [flags]
```

### Examples

```
  Detaches a public IP with ID "xxx" from a server with ID "yyy"
  $ stackit beta server public-ip detach xxx --server-id yyy
```

### Options

```
  -h, --help               Help for "stackit beta server public-ip detach"
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

* [stackit beta server public-ip](./stackit_beta_server_public-ip.md)	 - Allows attaching/detaching public IPs to servers

