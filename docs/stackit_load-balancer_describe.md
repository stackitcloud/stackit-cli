## stackit load-balancer describe

Shows details of a Load Balancer

### Synopsis

Shows details of a Load Balancer.

```
stackit load-balancer describe LOAD_BALANCER_NAME [flags]
```

### Examples

```
  Get details of a load balancer with name "my-load-balancer"
  $ stackit load-balancer describe my-load-balancer

  Get details of a load-balancer with name "my-load-balancer" in a JSON format
  $ stackit load-balancer describe my-load-balancer --output-format json
```

### Options

```
  -h, --help   Help for "stackit load-balancer describe"
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

* [stackit load-balancer](./stackit_load-balancer.md)	 - Provides functionality for Load Balancer

