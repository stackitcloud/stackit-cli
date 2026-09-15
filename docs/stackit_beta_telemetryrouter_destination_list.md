## stackit beta telemetryrouter destination list

Lists all TelemetryRouter destinations of an instance

### Synopsis

Lists all TelemetryRouter destinations of an instance.

```
stackit beta telemetryrouter destination list [flags]
```

### Examples

```
  Lists all destinations of the TelemetryRouter instance "xxx"
  $ stackit beta telemetryrouter destination list --instance-id xxx

  Lists all destinations in JSON format
  $ stackit beta telemetryrouter destination list --instance-id xxx --output-format json

  Lists up to 10 destinations
  $ stackit beta telemetryrouter destination list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit beta telemetryrouter destination list"
      --instance-id string   ID of the TelemetryRouter instance
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

* [stackit beta telemetryrouter destination](./stackit_beta_telemetryrouter_destination.md)	 - Provides functionality for TelemetryRouter destinations

