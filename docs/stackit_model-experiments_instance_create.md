## stackit model-experiments instance create

Creates an AI Model Experiments instance

### Synopsis

Creates an AI Model Experiments (MLflow) instance in your STACKIT project.

```
stackit model-experiments instance create [flags]
```

### Examples

```
  Create an AI Model Experiments instance with name "my-tracking"
  $ stackit model-experiments instance create --name my-tracking

  Create an instance with a description and labels
  $ stackit model-experiments instance create --name my-tracking --description "team tracking server" --labels env=prod,team=ml
```

### Options

```
      --deleted-experiment-retention string   Retention period for deleted experiments before permanent purge, e.g. "30d" (min 1d, max 90d)
      --description string                    Instance description
  -h, --help                                  Help for "stackit model-experiments instance create"
      --labels stringToString                 Labels as comma-separated key-value pairs, e.g. "--labels env=prod,team=ml" (default [])
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

* [stackit model-experiments instance](./stackit_model-experiments_instance.md)	 - Provides functionality for AI Model Experiments instances

