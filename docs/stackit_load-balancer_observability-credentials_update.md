## stackit load-balancer observability-credentials update

Updates observability credentials for Load Balancer

### Synopsis

Updates existing observability credentials (username and password) for Load Balancer. The credentials can be for Argus or another monitoring tool.

```
stackit load-balancer observability-credentials update [flags]
```

### Examples

```
  Update the password of observability credentials of Load Balancer with credentials reference "credentials-xxx". The password is entered using the terminal
  $ stackit load-balancer observability-credentials update credentials-xxx

  Update the password of observability credentials of Load Balancer with credentials reference "credentials-xxx", by providing the path to a file with the new password as flag
  $ stackit load-balancer observability-credentials update credentials-xxx --password @./new-password.txt
```

### Options

```
      --display-name string   Credentials name
  -h, --help                  Help for "stackit load-balancer observability-credentials update"
      --password string       Password. Can be a string or a file path, if prefixed with "@" (example: @./password.txt).
      --username string       Username
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit load-balancer observability-credentials](./stackit_load-balancer_observability-credentials.md)	 - Provides functionality for Load Balancer observability credentials

