## stackit beta alb-waf config list

Lists all ALB WAF configurations

### Synopsis

Lists all STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) configurations.

```
stackit beta alb-waf config list [flags]
```

### Examples

```
  List all ALB WAF configurations
  $ stackit beta alb-waf config list

  List the first 10 ALB WAF configurations
  $ stackit beta alb-waf config list --limit=10
```

### Options

```
  -h, --help        Help for "stackit beta alb-waf config list"
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

* [stackit beta alb-waf config](./stackit_beta_alb-waf_config.md)	 - Provides functionality for WAF configurations of the ALB WAF

