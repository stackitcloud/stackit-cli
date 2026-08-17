## stackit beta alb-waf custom-rule-group create

Creates an ALB WAF custom rule group

### Synopsis

Creates a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule group.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/alb-waf/version/v1 for information regarding the payload structure.

```
stackit beta alb-waf custom-rule-group create [flags]
```

### Examples

```
  Create an ALB WAF custom rule group using an API payload sourced from the file "./payload.json"
  $ stackit beta alb-waf custom-rule-group create --payload @./payload.json

  Create an ALB WAF custom rule group using an API payload provided as a JSON string
  $ stackit beta alb-waf custom-rule-group create --payload "{...}"
```

### Options

```
  -h, --help             Help for "stackit beta alb-waf custom-rule-group create"
      --payload string   Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).
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

