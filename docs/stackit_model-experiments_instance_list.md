## stackit model-experiments instance list

Lists all AI Model Experiments instances

### Synopsis

Lists all AI Model Experiments instances in a STACKIT project.

```
stackit model-experiments instance list [flags]
```

### Examples

```
  List AI Model Experiments instances in a project
  $ stackit model-experiments instance list

  List AI Model Experiments instances in a project in region eu02
  $ stackit model-experiments instance list --region eu02

  List AI Model Experiments instances in JSON format
  $ stackit model-experiments instance list --output-format json

  List up to 10 AI Model Experiments instances
  $ stackit model-experiments instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit model-experiments instance list"
      --limit int   Maximum number of entries to list
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

* [stackit model-experiments instance](./stackit_model-experiments_instance.md)	 - Provides functionality for AI Model Experiments instances

