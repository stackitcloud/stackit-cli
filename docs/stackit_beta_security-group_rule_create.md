## stackit beta security-group rule create

Creates a security group rule

### Synopsis

Creates a security group rule.

```
stackit beta security-group rule create [flags]
```

### Examples

```
  Create a security group rule for security group with ID "xxx" with direction "ingress"
  $ stackit beta security-group rule create --security-group-id xxx --direction ingress

  Create a security group rule for security group with ID "xxx" with direction "egress", protocol "icmp" and icmp parameters
  $ stackit beta security-group rule create --security-group-id xxx --direction egress --protocol-name icmp --icmp-parameter-code 0 --icmp-parameter-type 8

  Create a security group rule for security group with ID "xxx" with direction "ingress", protocol "tcp" and port range values
  $ stackit beta security-group rule create --security-group-id xxx --direction ingress --protocol-name tcp --port-range-max 24 --port-range-min 22

  Create a security group rule for security group with ID "xxx" with direction "ingress" and protocol number 1 
  $ stackit beta security-group rule create --security-group-id xxx --direction ingress --protocol-number 1
```

### Options

```
      --description string                The rule description
      --direction ingress                 The direction of the traffic which the rule should match. The possible values are: ingress, `egress`
      --ether-type string                 The ethertype which the rule should match
  -h, --help                              Help for "stackit beta security-group rule create"
      --icmp-parameter-code int           ICMP code. Can be set if the protocol is ICMP
      --icmp-parameter-type int           ICMP type. Can be set if the protocol is ICMP
      --ip-range string                   The remote IP range which the rule should match
      --port-range-max int                The maximum port number. Should be greater or equal to the minimum. This should only be provided if the protocol is not ICMP
      --port-range-min int                The minimum port number. Should be less or equal to the maximum. This should only be provided if the protocol is not ICMP
      --protocol-name protocol-name       The protocol name which the rule should match. If a protocol is to be defined, either protocol-name or `protocol-number` must be provided
      --protocol-number protocol-name     The protocol number which the rule should match. If a protocol is to be defined, either protocol-name or `protocol-number` must be provided
      --remote-security-group-id string   The remote security group which the rule should match
      --security-group-id string          The security group ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta security-group rule](./stackit_beta_security-group_rule.md)	 - Provides functionality for security group rules

