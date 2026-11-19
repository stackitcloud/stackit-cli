## stackit model-experiments token list

Lists all AI Model Experiments instance tokens

### Synopsis

Lists all auth tokens for an AI Model Experiments instance.

```
stackit model-experiments token list [flags]
```

### Examples

```
  List all tokens for an instance
  $ stackit model-experiments token list --instance-id xxx

  List tokens in region eu02
  $ stackit model-experiments token list --instance-id xxx --region eu02

  List tokens in JSON format
  $ stackit model-experiments token list --instance-id xxx --output-format json

  List up to 10 tokens for an instance
  $ stackit model-experiments token list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit model-experiments token list"
      --instance-id string   ID of the instance
      --limit int            Maximum number of entries to list
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

