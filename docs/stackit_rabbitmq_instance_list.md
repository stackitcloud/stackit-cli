## stackit rabbitmq instance list

Lists all RabbitMQ instances

### Synopsis

Lists all RabbitMQ instances.

```
stackit rabbitmq instance list [flags]
```

### Examples

```
  List all RabbitMQ instances
  $ stackit rabbitmq instance list

  List all RabbitMQ instances in JSON format
  $ stackit rabbitmq instance list --output-format json

  List up to 10 RabbitMQ instances
  $ stackit rabbitmq instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit rabbitmq instance list"
      --limit int   Maximum number of entries to list
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

* [stackit rabbitmq instance](./stackit_rabbitmq_instance.md)	 - Provides functionality for RabbitMQ instances

