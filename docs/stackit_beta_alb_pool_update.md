## stackit beta alb pool update

Updates an application target pool

### Synopsis

Updates an application target pool.

```
stackit beta alb pool update [flags]
```

### Examples

```
  Update an application target pool from a configuration file (the name of the pool is read from the file)
  $ stackit beta alb update --configuration my-target-pool.json --name my-load-balancer
```

### Options

```
  -c, --configuration string   filename of the input configuration file
  -h, --help                   Help for "stackit beta alb pool update"
  -n, --name string            name of the target pool name to update
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

* [stackit beta alb pool](./stackit_beta_alb_pool.md)	 - Manages target pools for application loadbalancers

