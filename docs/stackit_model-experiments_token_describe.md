## stackit model-experiments token describe

Shows details of an AI Model Experiments instance token

### Synopsis

Shows details of an auth token for an AI Model Experiments instance.

```
stackit model-experiments token describe TOKEN_ID [flags]
```

### Examples

```
  Show details of an auth token with ID "xxx"
  $ stackit model-experiments token describe xxx --instance-id yyy

  Show details of an auth token with ID "xxx" in JSON format
  $ stackit model-experiments token describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help[=true|false]    Help for "stackit model-experiments token describe"
      --instance-id string   ID of the instance
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit model-experiments token](./stackit_model-experiments_token.md)	 - Provides functionality for AI Model Experiments instance tokens

