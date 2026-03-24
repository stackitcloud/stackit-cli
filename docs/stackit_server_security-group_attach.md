## stackit server security-group attach

Attaches a security group to a server

### Synopsis

Attaches a security group to a server.

```
stackit server security-group attach [flags]
```

### Examples

```
  Attach a security group with ID "xxx" to a server with ID "yyy"
  $ stackit server security-group attach --server-id yyy --security-group-id xxx
```

### Options

```
  -h, --help                       Help for "stackit server security-group attach"
      --security-group-id string   Security Group ID
      --server-id string           Server ID
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

* [stackit server security-group](./stackit_server_security-group.md)	 - Allows attaching/detaching security groups to servers

