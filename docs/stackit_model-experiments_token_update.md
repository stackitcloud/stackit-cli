## stackit model-experiments token update

Updates an instance token

### Synopsis

Partially updates an auth token for an AI Model Experiments instance.

```
stackit model-experiments token update TOKEN_ID [flags]
```

### Examples

```
  Update an auth token with ID "xxx"
  $ stackit model-experiments token update xxx --instance-id yyy --name updated-token
```

### Options

```
      --description string      Token description
  -h, --help                    Help for "stackit model-experiments token update"
      --instance-id string      ID of the instance
      --labels stringToString   Labels as comma-separated key-value pairs, e.g. "--labels env=prod,team=ml" (default [])
      --name string             Token name
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

* [stackit model-experiments token](./stackit_model-experiments_token.md)	 - Provides functionality for AI Model Experiments instance tokens

