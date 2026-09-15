## stackit beta telemetryrouter instance delete

Deletes the given TelemetryRouter instance

### Synopsis

Deletes the given TelemetryRouter instance.

```
stackit beta telemetryrouter instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a TelemetryRouter instance with ID "xxx"
  $ stackit beta telemetryrouter instance delete "xxx"
```

### Options

```
  -h, --help   Help for "stackit beta telemetryrouter instance delete"
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

* [stackit beta telemetryrouter instance](./stackit_beta_telemetryrouter_instance.md)	 - Provides functionality for TelemetryRouter instances

