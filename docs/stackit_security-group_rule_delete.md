## stackit security-group rule delete

Deletes a security group rule

### Synopsis

Deletes a security group rule.
If the security group rule is still in use, the deletion will fail


```
stackit security-group rule delete SECURITY_GROUP_RULE_ID [flags]
```

### Examples

```
  Delete security group rule with ID "xxx" in security group with ID "yyy"
  $ stackit security-group rule delete xxx --security-group-id yyy
```

### Options

```
  -h, --help                       Help for "stackit security-group rule delete"
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

