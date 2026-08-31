## stackit beta alb-waf config describe

Shows details of an ALB WAF configuration

### Synopsis

Shows details of a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) configuration.

```
stackit beta alb-waf config describe NAME [flags]
```

### Examples

```
  Show details of an ALB WAF configuration with name "my-waf-config"
  $ stackit beta alb-waf config describe my-waf-config
```

### Options

```
  -h, --help   Help for "stackit beta alb-waf config describe"
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

