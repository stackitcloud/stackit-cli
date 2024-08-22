## stackit observability instance create

Creates an Observability instance

### Synopsis

Creates an Observability instance.

```
stackit observability instance create [flags]
```

### Examples

```
  Create an Observability instance with name "my-instance" and specify plan by name
  $ stackit observability instance create --name my-instance --plan-name Monitoring-Starter-EU01

  Create an Observability instance with name "my-instance" and specify plan by ID
  $ stackit observability instance create --name my-instance --plan-id xxx
```

### Options

```
  -h, --help               Help for "stackit observability instance create"
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
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit observability instance](./stackit_observability_instance.md)	 - Provides functionality for Observability instances

