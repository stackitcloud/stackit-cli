## stackit beta alb create

Creates an application loadbalancer

### Synopsis

Creates an application loadbalancer.

```
stackit beta alb create [flags]
```

### Examples

```
  Create an application loadbalancer from a configuration file
  $ stackit beta alb create --configuration my-loadbalancer.json
```

### Options

```
  -c, --configuration string   filename of the input configuration file
  -h, --help                   Help for "stackit beta alb create"
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

