## stackit model-experiments token create

Creates an AI Model Experiments instance token

### Synopsis

Creates an auth token for an AI Model Experiments instance.

```
stackit model-experiments token create [flags]
```

### Examples

```
  Create an auth token
  $ stackit model-experiments token create --instance-id xxx --name my-token

  Create an auth token with a description, labels, and time to live
  $ stackit model-experiments token create --instance-id xxx --name my-token --description "service token" --labels env=prod,team=ml --ttl-duration 30d
```

### Options

```
      --description string      Token description
  -h, --help                    Help for "stackit model-experiments token create"
      --instance-id string      ID of the instance
      --labels stringToString   Labels are comma-separated key-value string pairs, e.g. "--labels env=prod,team=ml" (default [])
      --name string             Token name
      --ttl-duration string     Maximum allowed duration is 6 months. If not set, defaults to 6 months. Format example: "720h" or "30d"
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

