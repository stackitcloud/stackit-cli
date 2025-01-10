## stackit rabbitmq plans

Lists all RabbitMQ service plans

### Synopsis

Lists all RabbitMQ service plans.

```
stackit rabbitmq plans [flags]
```

### Examples

```
  List all RabbitMQ service plans
  $ stackit rabbitmq plans

  List all RabbitMQ service plans in JSON format
  $ stackit rabbitmq plans --output-format json

  List up to 10 RabbitMQ service plans
  $ stackit rabbitmq plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit rabbitmq plans"
      --limit int   Maximum number of entries to list
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

* [stackit rabbitmq](./stackit_rabbitmq.md)	 - Provides functionality for RabbitMQ

