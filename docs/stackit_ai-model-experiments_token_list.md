## stackit ai-model-experiments token list

Lists instance tokens

### Synopsis

Lists all auth tokens for an AI Model Experiments instance.

```
stackit ai-model-experiments token list [flags]
```

### Examples

```
  List all tokens for an instance
  $ stackit ai-model-experiments token list --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit ai-model-experiments token list"
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

