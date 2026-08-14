## stackit beta edge-cloud instance describe

Describes an Edge Cloud instance

### Synopsis

Describes a STACKIT Edge Cloud (STEC) instance.

```
stackit beta edge-cloud instance describe INSTANCE_ID [flags]
```

### Examples

```
  Describe an Edge Cloud instance with ID "xxx"
  $ stackit beta edge-cloud instance describe xxx

  Get details of an Edge Cloud instance with ID "xxx" in JSON format
  $ stackit beta edge-cloud instance describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta edge-cloud instance describe"
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

