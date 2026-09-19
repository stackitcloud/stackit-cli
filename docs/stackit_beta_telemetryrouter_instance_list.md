## stackit beta telemetryrouter instance list

Lists TelemetryRouter instances

### Synopsis

Lists TelemetryRouter instances within the project.

```
stackit beta telemetryrouter instance list [flags]
```

### Examples

```
  List all TelemetryRouter instances
  $ stackit beta telemetryrouter instance list

  List the first 10 TelemetryRouter instances
  $ stackit beta telemetryrouter instance list --limit=10
```

### Options

```
  -h, --help        Help for "stackit beta telemetryrouter instance list"
      --limit int   Limit the output to the first n elements
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

