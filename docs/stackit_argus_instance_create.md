## stackit argus instance create

Creates an Argus instance

### Synopsis

Creates an Argus instance.

```
stackit argus instance create [flags]
```

### Examples

```
  Create an Argus instance with name "my-instance" and specify plan by name
  $ stackit argus instance create --name my-instance --plan-name Monitoring-Starter-EU01

  Create an Argus instance with name "my-instance" and specify plan by ID
  $ stackit argus instance create --name my-instance --plan-id xxx
```

### Options

```
  -h, --help               Help for "stackit argus instance create"
  -n, --name string        Instance name
      --plan-id string     Plan ID
      --plan-name string   Plan name
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit argus instance](./stackit_argus_instance.md)	 - Provides functionality for Argus instances

