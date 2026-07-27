## stackit beta edge-cloud instance update

Updates an Edge Cloud instance

### Synopsis

Updates a STACKIT Edge Cloud (STEC) instance.

```
stackit beta edge-cloud instance update INSTANCE_ID [flags]
```

### Examples

```
  Updates the description of an Edge Cloud instance with ID "xxx"
  $ stackit beta edge-cloud instance update xxx --description yyy

  Updates the plan of an Edge Cloud instance with ID "xxx"
  $ stackit beta edge-cloud instance update xxx --plan-id yyy
```

### Options

```
  -d, --description string   A user chosen description to distinguish multiple instances.
  -h, --help                 Help for "stackit beta edge-cloud instance update"
      --plan-id string       Service plan configures the size of the Instance.
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

* [stackit beta edge-cloud instance](./stackit_beta_edge-cloud_instance.md)	 - Provides functionality for Edge Cloud instances.

