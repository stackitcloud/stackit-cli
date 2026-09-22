## stackit ai-model-experiments instance patch

Updates an AI Model Experiments instance

### Synopsis

Partially updates an AI Model Experiments instance in a STACKIT project.

```
stackit ai-model-experiments instance patch INSTANCE_ID [flags]
```

### Examples

```
  Update the name of an AI Model Experiments instance with ID "xxx"
  $ stackit ai-model-experiments instance patch xxx --name my-new-name

  Update the description of an AI Model Experiments instance
  $ stackit ai-model-experiments instance patch xxx --description "team tracking server"

  Update labels on an AI Model Experiments instance
  $ stackit ai-model-experiments instance patch xxx --label env=prod

  Update multiple fields of an AI Model Experiments instance
  $ stackit ai-model-experiments instance patch xxx --name my-new-name --description "team tracking server" --label env=prod --deleted-experiment-retention 30d
```

### Options

```
      --deleted-experiment-retention string   Retention period for deleted experiments before permanent purge, e.g. "30d" (min 1d, max 90d)
      --description string                    Instance description
  -h, --help                                  Help for "stackit ai-model-experiments instance patch"
      --label stringToString                  Labels as key-value pairs, e.g. "--label env=prod" (default [])
  -n, --name string                           Instance name
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

* [stackit ai-model-experiments instance](./stackit_ai-model-experiments_instance.md)	 - Provides functionality for AI Model Experiments instances

