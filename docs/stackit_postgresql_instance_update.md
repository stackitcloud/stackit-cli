## stackit postgresql instance update

Updates a PostgreSQL instance

### Synopsis

Updates a PostgreSQL instance

```
stackit postgresql instance update [flags]
```

### Examples

```
$ stackit postgresql instance update --project-id xxx --instance-id xxx --plan-id xxx --acl xx.xx.xx.xx/xx
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            help for update
      --instance-id string              Instance ID
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
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

