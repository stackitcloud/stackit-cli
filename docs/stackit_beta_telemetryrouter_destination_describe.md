## stackit beta telemetryrouter destination describe

Shows details of a TelemetryRouter destination

### Synopsis

Shows details of a TelemetryRouter destination.

```
stackit beta telemetryrouter destination describe DESTINATION_ID [flags]
```

### Examples

```
  Show details of a TelemetryRouter destination with ID "xxx"
  $ stackit beta telemetryrouter destination describe xxx --instance-id yyy

  Show details of a TelemetryRouter destination with ID "xxx" in JSON format
  $ stackit beta telemetryrouter destination describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help                 Help for "stackit beta telemetryrouter destination describe"
      --instance-id string   ID of the TelemetryRouter instance
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

