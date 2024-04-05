## stackit rabbitmq credentials delete

Deletes credentials of a RabbitMQ instance

### Synopsis

Deletes credentials of a RabbitMQ instance.

```
stackit rabbitmq credentials delete CREDENTIALS_ID [flags]
```

### Examples

```
  Delete credentials with ID "xxx" of RabbitMQ instance with ID "yyy"
  $ stackit rabbitmq credentials delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit rabbitmq credentials delete"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit rabbitmq credentials](./stackit_rabbitmq_credentials.md)	 - Provides functionality for RabbitMQ credentials

