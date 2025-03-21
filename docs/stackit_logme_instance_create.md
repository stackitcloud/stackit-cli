## stackit logme instance create

Creates a LogMe instance

### Synopsis

Creates a LogMe instance.

```
stackit logme instance create [flags]
```

### Examples

```
  Create a LogMe instance with name "my-instance" and specify plan by name and version
  $ stackit logme instance create --name my-instance --plan-name stackit-logme2-1.2.50-replica --version 2

  Create a LogMe instance with name "my-instance" and specify plan by ID
  $ stackit logme instance create --name my-instance --plan-id xxx

  Create a LogMe instance with name "my-instance" and specify IP range which is allowed to access it
  $ stackit logme instance create --name my-instance --plan-id xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit logme instance create"
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
  -n, --name string                     Instance name
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --syslog strings                  Syslog
      --version string                  Instance LogMe version
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

* [stackit logme instance](./stackit_logme_instance.md)	 - Provides functionality for LogMe instances

