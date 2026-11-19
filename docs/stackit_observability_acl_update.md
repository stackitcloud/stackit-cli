## stackit observability acl update

Updates the ACLs of an Observability instance

### Synopsis

Updates the ACLs of an Observability instance.

```
stackit observability acl update [flags]
```

### Examples

```
  Set the ACLs of Observability instance "xxx" to the IP networks "1.2.3.0/24" and "5.6.7.8/32"
  $ stackit observability acl update --instance-id xxx --acl 1.2.3.0/24,5.6.7.8/32

  Remove all ACLs of Observability instance "xxx"
  $ stackit observability acl update --instance-id xxx --acl ""
```

### Options

```
      --acl strings          List of IP networks in CIDR notation which are allowed to access the instance. Replaces all existing ACLs and an empty value removes all ACLs (default [])
  -h, --help[=true|false]    Help for "stackit observability acl update"
      --instance-id string   Instance ID
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

