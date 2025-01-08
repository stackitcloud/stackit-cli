## stackit load-balancer quota

Shows the configured Load Balancer quota

### Synopsis

Shows the configured Load Balancer quota for the project. If you want to change the quota, please create a support ticket in the STACKIT Help Center (https://support.stackit.cloud)

```
stackit load-balancer quota [flags]
```

### Examples

```
  Get the configured load balancer quota for the project
  $ stackit load-balancer quota
```

### Options

```
  -h, --help   Help for "stackit load-balancer quota"
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

