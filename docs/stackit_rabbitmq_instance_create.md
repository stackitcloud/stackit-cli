## stackit rabbitmq instance create

Creates a RabbitMQ instance

### Synopsis

Creates a RabbitMQ instance.

```
stackit rabbitmq instance create [flags]
```

### Examples

```
  Create a RabbitMQ instance with name "my-instance" and specify plan by name and version
  $ stackit rabbitmq instance create --name my-instance --plan-name stackit-rabbitmq-1.2.10-replica --version 3.10

  Create a RabbitMQ instance with name "my-instance" and specify plan by ID
  $ stackit rabbitmq instance create --name my-instance --plan-id xxx

  Create a RabbitMQ instance with name "my-instance" and specify IP range which is allowed to access it
  $ stackit rabbitmq instance create --name my-instance --plan-id xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit rabbitmq instance create"
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
  -n, --name string                     Instance name
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --plugin strings                  Plugin
      --syslog strings                  Syslog
      --version string                  Instance RabbitMQ version
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

* [stackit rabbitmq instance](./stackit_rabbitmq_instance.md)	 - Provides functionality for RabbitMQ instances

