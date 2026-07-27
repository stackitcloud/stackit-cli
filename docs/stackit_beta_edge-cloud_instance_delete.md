## stackit beta edge-cloud instance delete

Deletes an Edge Cloud instance

### Synopsis

Deletes a STACKIT Edge Cloud (STEC) instance. The instance will be deleted permanently.

```
stackit beta edge-cloud instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete an edge instance with ID "xxx"
  $ stackit beta edge-cloud instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta edge-cloud instance delete"
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

