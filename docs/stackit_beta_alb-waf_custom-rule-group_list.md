## stackit beta alb-waf custom-rule-group list

Lists all ALB WAF custom rule groups

### Synopsis

Lists all STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule groups.

```
stackit beta alb-waf custom-rule-group list [flags]
```

### Examples

```
  List all ALB WAF custom rule groups
  $ stackit beta alb-waf custom-rule-group list

  List the first 10 ALB WAF custom rule groups
  $ stackit beta alb-waf custom-rule-group list --limit=10
```

### Options

```
  -h, --help        Help for "stackit beta alb-waf custom-rule-group list"
      --limit int   Limit the output to the first n elements
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

* [stackit beta alb-waf custom-rule-group](./stackit_beta_alb-waf_custom-rule-group.md)	 - Provides functionality for custom rule groups of the ALB WAF

