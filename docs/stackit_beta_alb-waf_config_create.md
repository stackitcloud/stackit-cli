## stackit beta alb-waf config create

Creates an ALB WAF configuration

### Synopsis

Creates a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) configuration.

```
stackit beta alb-waf config create [flags]
```

### Examples

```
  Create an ALB WAF configuration with name "my-waf-config"
  $ stackit beta alb-waf config create --name my-waf-config

  Create an ALB WAF configuration with a managed rule set, a custom rule group and labels
  $ stackit beta alb-waf config create --name my-waf-config --managed-rule-set-name my-managed-rule-set --custom-rule-group-name my-custom-rule-group --labels key1=value1,key2=value2
```

### Options

```
      --custom-rule-group-name string   Name of the custom rule group configuration to attach to the WAF
  -h, --help                            Help for "stackit beta alb-waf config create"
      --labels stringToString           Labels are key-value string pairs which can be attached to the WAF configuration. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --managed-rule-set-name string    Name of the managed rule set configuration to attach to the WAF
      --name string                     Name of the WAF configuration
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

* [stackit beta alb-waf config](./stackit_beta_alb-waf_config.md)	 - Provides functionality for WAF configurations of the ALB WAF

