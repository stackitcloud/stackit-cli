## stackit ai-model-experiments token patch

Updates an instance token

### Synopsis

Partially updates an auth token for an AI Model Experiments instance.

```
stackit ai-model-experiments token patch TOKEN_ID [flags]
```

### Examples

```
  Update an auth token with ID "xxx"
  $ stackit ai-model-experiments token patch xxx --instance-id yyy --name updated-token
```

### Options

```
      --description string     Token description
  -h, --help                   Help for "stackit ai-model-experiments token patch"
      --instance-id string     ID of the instance
      --label stringToString   Labels as key-value pairs, e.g. "--label env=prod" (default [])
      --name string            Token name
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

