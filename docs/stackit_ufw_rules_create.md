## stackit ufw rules create

Creates a UFW rule instance

### Synopsis

Creates a STACKIT Unified Firewall (UFW) rule instance.

```
stackit ufw rules create [flags]
```

### Examples

```
  Create a UFW rule instance of type ACL with sourceIp "1.1.1.1/32" of product "redis" for instance with id=ID
  $ stackit ufw instance create --product redis --sourceIp 1.1.1.1/32 --type ACL --instanceId ID
```

### Options

```
  -D, --description string       Description
  -d, --direction string         Direction (the direction of the traffic, typically ingress or egress, for security rules type)
  -e, --etherType string         Specifies the bound of the rule (for security rules type)
  -h, --help                     Help for "stackit ufw rules create"
  -i, --instanceId string        Instance ID that will have attached your rule
  -f, --offset int32             Offset - Position in the ACL list of an instance, will be ignored at creation (default -1)
  -r, --portRange string         Port range (the Port range to which the rule applies, for security rules type)
      --product string           The source service (e.g., Load Balancer, Redis) where you want to attach a rule
      --protocol string          The network protocol (e.g. TCP, UDP, ICMP, for security rules type)
  -g, --securityGroupId string   Security group ID - The ID of the Security Group
  -s, --sourceIp string          The IP (CIDR) to which the rule applies (e.g. 192.168.0.1/32)
  -t, --type string              Type (ACL/SecurityRule/SecurityGroup/PublicIP) You can check /provider-options route for them
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

