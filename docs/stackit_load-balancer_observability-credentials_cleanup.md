## stackit load-balancer observability-credentials cleanup

Deletes observability credentials unused by any Load Balancer

### Synopsis

Deletes observability credentials unused by any Load Balancer.

```
stackit load-balancer observability-credentials cleanup [flags]
```

### Examples

```
  Delete observability credentials unused by any Load Balancer
  $ stackit load-balancer observability-credentials cleanup
```

### Options

```
  -h, --help   Help for "stackit load-balancer observability-credentials cleanup"
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

* [stackit load-balancer observability-credentials](./stackit_load-balancer_observability-credentials.md)	 - Provides functionality for Load Balancer observability credentials

