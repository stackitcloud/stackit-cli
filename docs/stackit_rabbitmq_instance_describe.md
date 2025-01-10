## stackit rabbitmq instance describe

Shows details of a RabbitMQ instance

### Synopsis

Shows details of a RabbitMQ instance.

```
stackit rabbitmq instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a RabbitMQ instance with ID "xxx"
  $ stackit rabbitmq instance describe xxx

  Get details of a RabbitMQ instance with ID "xxx" in JSON format
  $ stackit rabbitmq instance describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit rabbitmq instance describe"
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

* [stackit rabbitmq instance](./stackit_rabbitmq_instance.md)	 - Provides functionality for RabbitMQ instances

