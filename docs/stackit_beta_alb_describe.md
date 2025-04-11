## stackit beta alb describe

Describes an application loadbalancer

### Synopsis

Describes an application loadbalancer.

```
stackit beta alb describe LOADBALANCER_NAME_ARG [flags]
```

### Examples

```
  Get details about an application loadbalancer with name "my-load-balancer"
  $ stackit beta alb describe my-load-balancer
```

### Options

```
  -h, --help   Help for "stackit beta alb describe"
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

* [stackit beta alb](./stackit_beta_alb.md)	 - Manages application loadbalancers

