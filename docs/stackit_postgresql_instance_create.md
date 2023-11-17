## stackit postgresql instance create

Creates a PostgreSQL instance

### Synopsis

Creates a PostgreSQL instance

```
stackit postgresql instance create [flags]
```

### Examples

```
$ stackit postgresql instance create --project-id xxx --name my-instance --plan-name plan-name --version version
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            help for create
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
  -n, --name string                     Instance name
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --plugin strings                  Plugin
      --syslog strings                  Syslog
      --version string                  Instance PostgreSQL version
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit postgresql instance](./stackit_postgresql_instance.md)	 - Provides functionality for PostgreSQL instance

