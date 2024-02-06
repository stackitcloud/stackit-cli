## stackit mariadb instance create

Creates a MariaDB instance

### Synopsis

Creates a MariaDB instance.

```
stackit mariadb instance create [flags]
```

### Examples

```
  Create a MariaDB instance with name "my-instance" and specify plan by name and version
  $ stackit mariadb instance create --name my-instance --plan-name stackit-mariadb-1.2.10-replica --version 10.6

  Create a MariaDB instance with name "my-instance" and specify plan by ID
  $ stackit mariadb instance create --name my-instance --plan-id xxx

  Create a MariaDB instance with name "my-instance" and specify IP range which is allowed to access it
  $ stackit mariadb instance create --name my-instance --plan-id xxx --acl 192.168.1.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit mariadb instance create"
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
  -n, --name string                     Instance name
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --plugin strings                  Plugin
      --syslog strings                  Syslog
      --version string                  Instance MariaDB version
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit mariadb instance](./stackit_mariadb_instance.md)	 - Provides functionality for MariaDB instances

