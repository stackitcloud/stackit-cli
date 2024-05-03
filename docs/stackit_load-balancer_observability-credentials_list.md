## stackit load-balancer observability-credentials list

Lists all observability credentials for Load Balancer

### Synopsis

Lists all observability credentials for Load Balancer.

```
stackit load-balancer observability-credentials list [flags]
```

### Examples

```
  List all observability credentials for Load Balancer
  $ stackit load-balancer observability-credentials list

  List all observability credentials for Load Balancer in JSON format
  $ stackit load-balancer observability-credentials list --output-format json

  List up to 10 observability credentials for Load Balancer
  $ stackit load-balancer observability-credentials list --limit 10
```

### Options

```
  -h, --help        Help for "stackit load-balancer observability-credentials list"
      --limit int   Maximum number of entries to list
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

