## stackit beta alb-waf managed-rule-set create

Creates a managed rule set for the ALB WAF

### Synopsis

Creates a managed rule set (MRS) for the Web Application Firewall (WAF) of application loadbalancers.

```
stackit beta alb-waf managed-rule-set create [flags]
```

### Examples

```
  Create a managed rule set with name "my-managed-rule-set"
  $ stackit beta alb-waf managed-rule-set create --name my-managed-rule-set
```

### Options

```
  -h, --help          Help for "stackit beta alb-waf managed-rule-set create"
      --name string   Name of the managed rule set
      --type string   Managed rule set type (one of: [TYPE_OWASP_CRS]) (default "TYPE_OWASP_CRS")
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

