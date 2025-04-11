## stackit beta alb template

creates configuration templates to use for resource creation

### Synopsis

creates a json or yaml template file for creating/updating an application loadbalancer or target pool.

```
stackit beta alb template [flags]
```

### Examples

```
  Create a yaml template
  $ stackit beta alb template --format=yaml --type alb

  Create a json template
  $ stackit beta alb template --format=json --type pool
```

### Options

```
  -f, --format string   Defines the output format ('yaml' or 'json'), default is 'json' (default "json")
  -h, --help            Help for "stackit beta alb template"
  -t, --type string     Defines the output type ('alb' or 'pool'), default is 'alb' (default "alb")
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

