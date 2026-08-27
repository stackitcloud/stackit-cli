## stackit valkey instance update

Updates a Valkey instance

### Synopsis

Updates a Valkey instance.

```
stackit valkey instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the plan of a Valkey instance with ID "xxx" by plan ID
  $ stackit valkey instance update xxx --plan-id yyy

  Update the plan of a Valkey instance with ID "xxx" by name and version
  $ stackit valkey instance update xxx --plan-name stackit-keyvalue-1.2.10-replica --version 8

  Update the range of IPs allowed to access a Valkey instance with ID "xxx"
  $ stackit valkey instance update xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit valkey instance update"
      --metrics-frequency int32         Metrics frequency in seconds
      --metrics-prefix string           Metrics prefix
      --min-replicas-to-write int32     Minimum number of replicas that must acknowledge a write for it to be accepted
      --monitoring-instance-id string   Monitoring instance ID
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --repl-backlog-size string        Replication backlog size (e.g. "1mb")
      --syslog strings                  Syslog
      --version string                  Instance Valkey version
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit valkey instance](./stackit_valkey_instance.md)	 - Provides functionality for Valkey instances

