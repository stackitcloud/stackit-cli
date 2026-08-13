## stackit beta alb-waf custom-rule-group generate-payload

Generates a payload to create/update an ALB WAF custom rule group

### Synopsis

Generates a JSON payload with values to be used as --payload input for ALB WAF custom rule group creation or update.
If --name is set, the payload is generated with the current values of the given custom rule group, to be used with the update command. If unset, a payload with default values for the create command is generated.
See https://docs.api.stackit.cloud/documentation/alb-waf/version/v1 for information regarding the payload structure.

```
stackit beta alb-waf custom-rule-group generate-payload [flags]
```

### Examples

```
  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit beta alb-waf custom-rule-group generate-payload --file-path ./payload.json
  <Modify payload in file, if needed>
  $ stackit beta alb-waf custom-rule-group create --payload @./payload.json

  Generate a payload with the current values of an existing custom rule group, and adapt it with custom values for the different configuration options
  $ stackit beta alb-waf custom-rule-group generate-payload --name my-custom-rule-group --file-path ./payload.json
  <Modify payload in file>
  $ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload @./payload.json

  Generate a payload with the current values of an existing custom rule group, and preview it in the terminal
  $ stackit beta alb-waf custom-rule-group generate-payload --name my-custom-rule-group
```

### Options

```
  -f, --file-path string   If set, writes the payload to the given file. If unset, writes the payload to the standard output
  -h, --help               Help for "stackit beta alb-waf custom-rule-group generate-payload"
  -n, --name string        If set, generates the payload with the current values of the given custom rule group (to be used with the update command). If unset, generates the payload with default values (to be used with the create command)
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

