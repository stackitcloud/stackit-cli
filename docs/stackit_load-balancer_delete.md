## stackit load-balancer delete

Deletes a Load Balancer

### Synopsis

Deletes a Load Balancer.

```
stackit load-balancer delete LOAD_BALANCER_NAME [flags]
```

### Examples

```
  Deletes a load balancer with name "my-load-balancer"
  $ stackit load-balancer delete my-load-balancer
```

### Options

```
  -h, --help   Help for "stackit load-balancer delete"
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

