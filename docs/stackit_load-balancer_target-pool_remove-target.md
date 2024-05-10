## stackit load-balancer target-pool remove-target

Removes a target from a target pool

### Synopsis

Removes a target from a target pool.

```
stackit load-balancer target-pool remove-target TARGET_IP [flags]
```

### Examples

```
  Remove target with IP 1.2.3.4 from target pool "my-target-pool" of load balancer with name "my-load-balancer"
  $ stackit load-balancer target-pool remove-target 1.2.3.4 --target-pool-name my-target-pool --lb-name my-load-balancer
```

### Options

```
  -h, --help                      Help for "stackit load-balancer target-pool remove-target"
      --lb-name string            Load balancer name
      --target-pool-name string   Target IP of the target to remove. Must be a valid IPv4 or IPv6
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

