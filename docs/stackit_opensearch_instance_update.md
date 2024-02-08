## stackit opensearch instance update

Updates an OpenSearch instance

### Synopsis

Updates an OpenSearch instance.

```
stackit opensearch instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the plan of an OpenSearch instance with ID "xxx"
  $ stackit opensearch instance update xxx --plan-id yyy

  Update the range of IPs allowed to access an OpenSearch instance with ID "xxx"
  $ stackit opensearch instance update xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings                     List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --enable-monitoring               Enable monitoring
      --graphite string                 Graphite host
  -h, --help                            Help for "stackit opensearch instance update"
      --metrics-frequency int           Metrics frequency
      --metrics-prefix string           Metrics prefix
      --monitoring-instance-id string   Monitoring instance ID
      --plan-id string                  Plan ID
      --plan-name string                Plan name
      --plugin strings                  Plugin
      --syslog strings                  Syslog
      --version string                  Instance OpenSearch version
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit opensearch instance](./stackit_opensearch_instance.md)	 - Provides functionality for OpenSearch instances

