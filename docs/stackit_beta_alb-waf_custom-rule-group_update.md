## stackit beta alb-waf custom-rule-group update

Updates an ALB WAF custom rule group

### Synopsis

Updates a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule group.
The rules of the custom rule group are replaced atomically: send the complete desired rule set, per-rule partial merge is not supported.
The payload can be provided as a JSON string or a file path prefixed with "@". The rule IDs and the rule behavior severity are managed by the server and cannot be set.
See https://docs.api.stackit.cloud/documentation/alb-waf/version/v1 for information regarding the payload structure.

```
stackit beta alb-waf custom-rule-group update CUSTOM_RULE_GROUP_NAME [flags]
```

### Examples

```
  Update an ALB WAF custom rule group using an API payload sourced from the file "./payload.json"
  $ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload @./payload.json

  Update an ALB WAF custom rule group using an API payload provided as a JSON string
  $ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload "{...}"

  Generate a payload with the current values of an existing custom rule group, adapt it and update the custom rule group with it
  $ stackit beta alb-waf custom-rule-group generate-payload --name my-custom-rule-group > ./payload.json
  <Modify payload in file>
  $ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload @./payload.json
```

### Options

```
  -h, --help             Help for "stackit beta alb-waf custom-rule-group update"
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

