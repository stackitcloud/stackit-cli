## stackit load-balancer target-pool add-target

Adds a target to a target pool

### Synopsis

Adds a target to a target pool.

```
stackit load-balancer target-pool add-target LOAD_BALANCER_NAME [flags]
```

### Examples

```
  Add a target to target pool "my-target-pool" of load balancer with name "my-load-balancer"
  $ stackit load-balancer target-pool add-target my-load-balancer --target-pool-name my-target-pool --target-name my-new-target --ip 1.2.3.4
```

### Options

```
  -h, --help                      Help for "stackit load-balancer target-pool add-target"
      --ip string                 Target IP. Must by unique within a target pool. Must be a valid IPv4 or IPv6
      --target-name string        Target name
      --target-pool-name string   Target pool name
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

* [stackit load-balancer target-pool](./stackit_load-balancer_target-pool.md)	 - Provides functionality for target pools

