## stackit ai-model-experiments token delete

Deletes an instance token

### Synopsis

Deletes an auth token from an AI Model Experiments instance.

```
stackit ai-model-experiments token delete TOKEN_ID [flags]
```

### Examples

```
  Delete an auth token with ID "xxx"
  $ stackit ai-model-experiments token delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit ai-model-experiments token delete"
      --instance-id string   ID of the instance
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

* [stackit ai-model-experiments token](./stackit_ai-model-experiments_token.md)	 - Provides functionality for AI Model Experiments instance tokens

