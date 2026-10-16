## stackit ufw rules update

Updates a UFW rule instance

### Synopsis

Updates a STACKIT Unified Firewall (UFW) rule instance.

```
stackit ufw rules update INSTANCE_ID [flags]
```

### Examples

```
  Update a UFW rule instance with "1.1.1.1/32" as sourceIp for instance with ID "xxx"
  $ stackit ufw rules update xxx --sourceIp 1.1.1.1/32
```

### Options

```
  -h, --help              Help for "stackit ufw rules update"
  -s, --sourceIp string   The IP (CIDR) to which the rule applies (e.g. 192.168.0.1/32)
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

* [stackit ufw rules](./stackit_ufw_rules.md)	 - Provides functionality for UFW rules

