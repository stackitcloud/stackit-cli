## stackit redis instance create

Creates a Redis instance

### Synopsis

Creates a Redis instance.

```
stackit redis instance create [flags]
```

### Examples

```
  Create a Redis instance with name "my-instance" and specify plan by name and version
  $ stackit redis instance create --name my-instance --plan-name stackit-redis-1.2.10-replica --version 6

  Create a Redis instance with name "my-instance" and specify plan by ID
  $ stackit redis instance create --name my-instance --plan-id xxx

  Create a Redis instance with name "my-instance" and specify IP range which is allowed to access it
  $ stackit redis instance create --name my-instance --plan-id xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit redis instance create"
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
  -n, --name string                     Instance name
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --syslog strings                  Syslog
      --version string                  Instance Redis version
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit redis instance](./stackit_redis_instance.md)	 - Provides functionality for Redis instances

