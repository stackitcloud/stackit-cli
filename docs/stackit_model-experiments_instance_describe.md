## stackit model-experiments instance describe

Shows details of an AI Model Experiments instance

### Synopsis

Shows details of an AI Model Experiments instance in a STACKIT project.

```
stackit model-experiments instance describe INSTANCE_ID [flags]
```

### Examples

```
  Show details of an AI Model Experiments instance with ID "xxx"
  $ stackit model-experiments instance describe xxx
```

### Options

```
  -h, --help   Help for "stackit model-experiments instance describe"
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

