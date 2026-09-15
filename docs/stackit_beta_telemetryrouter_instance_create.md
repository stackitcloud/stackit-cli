## stackit beta telemetryrouter instance create

Creates a TelemetryRouter instance

### Synopsis

Creates a TelemetryRouter instance.

```
stackit beta telemetryrouter instance create [flags]
```

### Examples

```
  Create a TelemetryRouter instance with name "my-instance"
  $ stackit beta telemetryrouter instance create --display-name "my-instance"

  Create a TelemetryRouter instance with name "my-instance" and a description
  $ stackit beta telemetryrouter instance create --display-name "my-instance" --description "Description of the instance"

  Create a TelemetryRouter instance with name "my-instance" that only routes telemetry data for HTTP GET or HEAD requests, matched at the resource level
  $ stackit beta telemetryrouter instance create --display-name "my-instance" --filter-attribute "key=http.method;level=resource;matcher==;values=GET,HEAD"

  Create a TelemetryRouter instance with name "my-instance" and many filter attributes sourced from a JSON file (recommended over repeating --filter-attribute for more than a handful of attributes); use "stackit beta telemetryrouter instance generate-filter" to generate a starting point for "./filter.json"
  $ stackit beta telemetryrouter instance create --display-name "my-instance" --filter @./filter.json
```

### Options

```
      --description string             Description
      --display-name string            Display name
      --filter string                  Filter configuration as JSON, of the form {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...]}. Can be a JSON string or a file path, if prefixed with "@" (example: @./filter.json). Recommended over "--filter-attribute" when setting many filter attributes at once. Run "stackit beta telemetryrouter instance generate-filter" to generate a starting point. Mutually exclusive with "--filter-attribute".
      --filter-attribute stringArray   Adds a filter attribute that restricts which telemetry data is routed. Must be of the form "key=<attribute-key>;level=<resource|scope|logRecord>;matcher=<=|!=>;values=<v1>,<v2>,...", e.g. "key=http.method;level=resource;matcher==;values=GET,HEAD" only routes data whose "http.method" resource attribute equals "GET" or "HEAD". Can be repeated to add multiple filter attributes, which are combined with AND. Impractical for many filter attributes at once; use "--filter" instead. Mutually exclusive with "--filter".
  -h, --help                           Help for "stackit beta telemetryrouter instance create"
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

