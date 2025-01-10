## stackit load-balancer list

Lists all Load Balancers

### Synopsis

Lists all Load Balancers.

```
stackit load-balancer list [flags]
```

### Examples

```
  List all load balancers
  $ stackit load-balancer list

  List all loadbalancers in JSON format
  $ stackit load-balancer list --output-format json

  List up to 10 load balancers 
  $ stackit load-balancer list --limit 10
```

### Options

```
  -h, --help        Help for "stackit load-balancer list"
      --limit int   Maximum number of entries to list
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

* [stackit load-balancer](./stackit_load-balancer.md)	 - Provides functionality for Load Balancer

