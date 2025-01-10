## stackit observability instance update

Updates an Observability instance

### Synopsis

Updates an Observability instance.

```
stackit observability instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the plan of an Observability instance with ID "xxx" by specifying the plan ID
  $ stackit observability instance update xxx --plan-id yyy

  Update the plan of an Observability instance with ID "xxx" by specifying the plan name
  $ stackit observability instance update xxx --plan-name Frontend-Starter-EU01

  Update the name of an Observability instance with ID "xxx"
  $ stackit observability instance update xxx --name new-instance-name
```

### Options

```
  -h, --help               Help for "stackit observability instance update"
  -n, --name string        Instance name
      --plan-id string     Plan ID
      --plan-name string   Plan name
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

* [stackit observability instance](./stackit_observability_instance.md)	 - Provides functionality for Observability instances

