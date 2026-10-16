## stackit ufw rules create

Creates a UFW rule instance

### Synopsis

Creates a STACKIT Unified Firewall (UFW) rule instance.

```
stackit ufw rules create [flags]
```

### Examples

```
  Create a UFW rule instance of type ACL with sourceIp "1.1.1.1/32" of product "Redis" for instance with id=ID
  $ stackit ufw rules create --product redis --sourceIp 1.1.1.1/32 --type ACL --instanceId ID

  Create a UFW rule instance of type ACL with sourceIp "2.2.2.2/32" of product "Edge Cloud" for instance with id=ID
  $ stackit ufw rules create --product edge-cloud --sourceIp 2.2.2.2/32 --type ACL --instanceId ID
```

### Options

```
  -h, --help                Help for "stackit ufw rules create"
  -i, --instanceId string   Instance ID that will have attached your rule
      --product string      The source service (e.g., Edge Cloud, Redis) where you want to attach a rule
  -s, --sourceIp string     The IP (CIDR) to which the rule applies (e.g. 192.168.0.1/32)
  -t, --type string         Type (ACL/SecurityRule/SecurityGroup) You can check /provider-options route for them. Unfortunately, this field could be only ACL for the CLI version
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

