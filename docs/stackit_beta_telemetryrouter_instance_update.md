## stackit beta telemetryrouter instance update

Updates a TelemetryRouter instance

### Synopsis

Updates a TelemetryRouter instance.

```
stackit beta telemetryrouter instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the display name of the TelemetryRouter instance with ID "xxx"
  $ stackit beta telemetryrouter instance update xxx --display-name new-name

  Update the description of the TelemetryRouter instance with ID "xxx"
  $ stackit beta telemetryrouter instance update xxx --description new-description

  Replace the filter attributes of the TelemetryRouter instance with ID "xxx" so it only routes telemetry data for HTTP GET or HEAD requests, matched at the resource level
  $ stackit beta telemetryrouter instance update xxx --filter-attribute "key=http.method;level=resource;matcher==;values=GET,HEAD"

  Replace the filter attributes of the TelemetryRouter instance with ID "xxx" with many filter attributes sourced from a JSON file (recommended over repeating --filter-attribute for more than a handful of attributes); use "stackit beta telemetryrouter instance generate-filter --instance-id xxx" to generate a starting point from the instance's current filter
  $ stackit beta telemetryrouter instance update xxx --filter @./filter.json
```

### Options

```
      --description string             Description
      --display-name string            Display name
      --filter string                  Filter configuration as JSON, of the form {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...]}. Can be a JSON string or a file path, if prefixed with "@" (example: @./filter.json). Replaces any existing filter attributes on the instance. Recommended over "--filter-attribute" when setting many filter attributes at once. Run "stackit beta telemetryrouter instance generate-filter --instance-id <id>" to generate a starting point from the instance's current filter. Mutually exclusive with "--filter-attribute". The API does not support clearing all filter attributes on update; at least one filter attribute must remain set.
      --filter-attribute stringArray   Adds a filter attribute that restricts which telemetry data is routed. Must be of the form "key=<attribute-key>;level=<resource|scope|logRecord>;matcher=<=|!=>;values=<v1>,<v2>,...", e.g. "key=http.method;level=resource;matcher==;values=GET,HEAD" only routes data whose "http.method" resource attribute equals "GET" or "HEAD". Can be repeated to add multiple filter attributes, which are combined with AND. Replaces any existing filter attributes on the instance. Impractical for many filter attributes at once; use "--filter" instead. Mutually exclusive with "--filter". The API does not support clearing all filter attributes on update; at least one filter attribute must remain set.
  -h, --help                           Help for "stackit beta telemetryrouter instance update"
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

