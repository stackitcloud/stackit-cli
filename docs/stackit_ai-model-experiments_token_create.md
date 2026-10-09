## stackit ai-model-experiments token create

Creates an instance token

### Synopsis

Creates an auth token for an AI Model Experiments instance.

```
stackit ai-model-experiments token create [flags]
```

### Examples

```
  Create an auth token
  $ stackit ai-model-experiments token create --instance-id xxx --name my-token
```

### Options

```
      --description string     Token description
  -h, --help                   Help for "stackit ai-model-experiments token create"
      --instance-id string     ID of the instance
      --label stringToString   Labels as key-value pairs, e.g. "--label env=prod" (default [])
      --name string            Token name
      --ttl-duration string    Token time to live duration
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

