## stackit model-experiments token update

Updates an AI Model Experiments instance token

### Synopsis

Partially updates an auth token for an AI Model Experiments instance.

```
stackit model-experiments token update TOKEN_ID [flags]
```

### Examples

```
  Update the name of an auth token with ID "xxx"
  $ stackit model-experiments token update xxx --instance-id yyy --name updated-token

  Update the description of an auth token
  $ stackit model-experiments token update xxx --instance-id yyy --description "updated token"

  Update labels on an auth token
  $ stackit model-experiments token update xxx --instance-id yyy --labels env=prod,team=ml

  Update multiple fields of an auth token
  $ stackit model-experiments token update xxx --instance-id yyy --name updated-token --description "updated token" --labels env=prod,team=ml
```

### Options

```
      --description string      Token description
  -h, --help                    Help for "stackit model-experiments token update"
      --instance-id string      ID of the instance
      --labels stringToString   Labels are comma-separated key-value string pairs, e.g. "--labels env=prod,team=ml" (default [])
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

