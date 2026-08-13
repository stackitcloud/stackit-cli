## stackit beta alb-waf custom-rule-group delete

Deletes an ALB WAF custom rule group

### Synopsis

Deletes a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule group.
A custom rule group can only be deleted if it is not referenced by any WAF configuration.

```
stackit beta alb-waf custom-rule-group delete CUSTOM_RULE_GROUP_NAME [flags]
```

### Examples

```
  Delete an ALB WAF custom rule group with name "my-custom-rule-group"
  $ stackit beta alb-waf custom-rule-group delete my-custom-rule-group
```

### Options

```
  -h, --help   Help for "stackit beta alb-waf custom-rule-group delete"
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

* [stackit beta alb-waf custom-rule-group](./stackit_beta_alb-waf_custom-rule-group.md)	 - Provides functionality for alb-waf Custom Rule Group

