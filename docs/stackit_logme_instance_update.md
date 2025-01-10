## stackit logme instance update

Updates a LogMe instance

### Synopsis

Updates a LogMe instance.

```
stackit logme instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the plan of a LogMe instance with ID "xxx"
  $ stackit logme instance update xxx --plan-id yyy

  Update the range of IPs allowed to access a LogMe instance with ID "xxx"
  $ stackit logme instance update xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit logme instance update"
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
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

