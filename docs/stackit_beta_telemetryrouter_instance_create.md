## stackit beta telemetryrouter instance create

Creates a Telemetry Router instance

### Synopsis

Creates a Telemetry Router instance.

```
stackit beta telemetryrouter instance create [flags]
```

### Examples

```
  Create a Telemetry Router instance with name "my-instance"
  $ stackit beta telemetryrouter instance create --display-name "my-instance"

  Create a Telemetry Router instance with name "my-instance" and a description
  $ stackit beta telemetryrouter instance create --display-name "my-instance" --description "Description of the instance"

  Create a Telemetry Router instance with name "my-instance" that only routes telemetry data for HTTP GET or HEAD requests, matched at the resource level, using a filter configuration sourced from a JSON file; use "stackit beta telemetryrouter instance generate-filter" to generate a starting point for "./filter.json"; available attributes are documented in https://docs.stackit.cloud/products/logging-and-monitoring/telemetry-router/basics/dataformats/
  $ stackit beta telemetryrouter instance create --display-name "my-instance" --filter @./filter.json
```

### Options

```
      --description string    Description
      --display-name string   Display name
      --filter string         Filter configuration as JSON, of the form {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...]}. Can be a JSON string or a file path, if prefixed with "@" (example: @./filter.json). Run "stackit beta telemetryrouter instance generate-filter" to generate a starting point.
  -h, --help                  Help for "stackit beta telemetryrouter instance create"
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

* [stackit beta telemetryrouter instance](./stackit_beta_telemetryrouter_instance.md)	 - Provides functionality for Telemetry Router instances

