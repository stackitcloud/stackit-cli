## stackit load-balancer observability-credentials describe

Shows details of observability credentials for Load Balancer

### Synopsis

Shows details of observability credentials for Load Balancer.

```
stackit load-balancer observability-credentials describe CREDENTIALS_REF [flags]
```

### Examples

```
  Get details of observability credentials with reference "credentials-xxx"
  $ stackit load-balancer observability-credentials describe credentials-xxx
```

### Options

```
  -h, --help   Help for "stackit load-balancer observability-credentials describe"
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

