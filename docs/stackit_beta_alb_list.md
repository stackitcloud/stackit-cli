## stackit beta alb list

Lists albs

### Synopsis

Lists application load balancers.

```
stackit beta alb list [flags]
```

### Examples

```
  List all load balancers
  $ stackit beta alb list

  List the first 10 application load balancers
  $ stackit beta alb list --limit=10
```

### Options

```
  -h, --help        Help for "stackit beta alb list"
      --limit int   Limit the output to the first n elements
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

