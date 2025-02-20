## stackit security-group rule describe

Shows details of a security group rule

### Synopsis

Shows details of a security group rule.

```
stackit security-group rule describe SECURITY_GROUP_RULE_ID [flags]
```

### Examples

```
  Show details of a security group rule with ID "xxx" in security group with ID "yyy"
  $ stackit security-group rule describe xxx --security-group-id yyy

  Show details of a security group rule with ID "xxx" in security group with ID "yyy" in JSON format
  $ stackit security-group rule describe xxx --security-group-id yyy --output-format json
```

### Options

```
  -h, --help                       Help for "stackit security-group rule describe"
      --security-group-id string   The security group ID
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

* [stackit security-group rule](./stackit_security-group_rule.md)	 - Provides functionality for security group rules

