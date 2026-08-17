## stackit beta alb-waf managed-rule-set describe

Describes a managed rule set of the ALB WAF

### Synopsis

Describes a managed rule set (MRS) of the Web Application Firewall (WAF) for application loadbalancers.

```
stackit beta alb-waf managed-rule-set describe NAME [flags]
```

### Examples

```
  Get details about a managed rule set with name "my-managed-rule-set"
  $ stackit beta alb-waf managed-rule-set describe my-managed-rule-set
```

### Options

```
  -h, --help   Help for "stackit beta alb-waf managed-rule-set describe"
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

* [stackit beta alb-waf managed-rule-set](./stackit_beta_alb-waf_managed-rule-set.md)	 - Provides functionality for managed rule sets of the ALB WAF

