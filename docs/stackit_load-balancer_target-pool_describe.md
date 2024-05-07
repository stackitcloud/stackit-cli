## stackit load-balancer target-pool describe

Shows details of a target pool in a Load Balancer

### Synopsis

Shows details of a target pool in a Load Balancer.

```
stackit load-balancer target-pool describe TARGET_POOL_NAME [flags]
```

### Examples

```
  Get details of a target pool with name "pool" in load balancer with name "my-load-balancer"
  $ stackit load-balancer target-pool describe pool --lb-name my-load-balancer

  Get details of a target pool with name "pool" in load balancer with name "my-load-balancer in JSON output"
  $ stackit load-balancer target-pool describe pool --lb-name my-load-balancer --output-format json
```

### Options

```
  -h, --help             Help for "stackit load-balancer target-pool describe"
      --lb-name string   Name of the load balancer
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

