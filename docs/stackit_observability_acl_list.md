## stackit observability acl list

Lists all ACLs of an Observability instance

### Synopsis

Lists all IP networks in CIDR notation which are allowed to access an Observability instance.

```
stackit observability acl list [flags]
```

### Examples

```
  List all ACLs of Observability instance "xxx"
  $ stackit observability acl list --instance-id xxx

  List all ACLs of Observability instance "xxx" in JSON format
  $ stackit observability acl list --instance-id xxx --output-format json

  List up to 10 ACLs of Observability instance "xxx"
  $ stackit observability acl list --instance-id xxx --limit 10
```

### Options

```
  -h, --help[=true|false]    Help for "stackit observability acl list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit observability acl](./stackit_observability_acl.md)	 - Provides functionality for ACLs of Observability instances

