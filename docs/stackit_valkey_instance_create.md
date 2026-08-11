## stackit valkey instance create

Creates a Valkey instance

### Synopsis

Creates a Valkey instance.

```
stackit valkey instance create [flags]
```

### Examples

```
  Create a Valkey instance with name "my-instance" and specify plan by name and version
  $ stackit beta valkey instance create --name my-instance --plan-name stackit-keyvalue-1.2.10-replica --version 8

  Create a Valkey instance with name "my-instance" and specify plan by ID
  $ stackit beta valkey instance create --name my-instance --plan-id xxx

  Create a Valkey instance with name "my-instance" and specify IP range which is allowed to access it
  $ stackit beta valkey instance create --name my-instance --plan-id xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit valkey instance create"
      --metrics-frequency int32         Metrics frequency in seconds
      --metrics-prefix string           Metrics prefix
      --min-replicas-to-write int32     Minimum number of replicas that must acknowledge a write for it to be accepted (Valkey only)
      --monitoring-instance-id string   Monitoring instance ID
  -n, --name string                     Instance name
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --repl-backlog-size string        Replication backlog size (e.g. "1mb") (Valkey only)
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

