## stackit beta alb-waf managed-rule-set update

Updates a managed rule set of the ALB WAF

### Synopsis

Updates the rules of a managed rule set (MRS) of the Web Application Firewall (WAF) for application loadbalancers. Only the rules provided in the configuration file are updated, all other rules remain unchanged.

```
stackit beta alb-waf managed-rule-set update NAME [flags]
```

### Examples

```
  Update the rules of a managed rule set with name "my-managed-rule-set" from a configuration file
  $ stackit beta alb-waf managed-rule-set update my-managed-rule-set --configuration my-rules.json
```

### Options

```
  -c, --configuration string   Filename of the input configuration file
  -h, --help                   Help for "stackit beta alb-waf managed-rule-set update"
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

