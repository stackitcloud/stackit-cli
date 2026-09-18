## stackit beta telemetryrouter destination create

Creates a TelemetryRouter destination

### Synopsis

Creates a destination for a TelemetryRouter instance.
The payload can be provided as a JSON string or a file path prefixed with "@".
If no payload is provided, a default payload will be used.
The "config" field is a discriminated union: exactly one of "openTelemetry" or "s3" must be set, matching the "configType".

```
stackit beta telemetryrouter destination create [flags]
```

### Examples

```
  Create a destination on TelemetryRouter instance "xxx" using the default (OpenTelemetry) configuration
  $ stackit beta telemetryrouter destination create --instance-id xxx

  Create a destination on TelemetryRouter instance "xxx" using an API payload sourced from the file "./payload.json"
  $ stackit beta telemetryrouter destination create --payload @./payload.json --instance-id xxx

  Create a destination on TelemetryRouter instance "xxx" using an API payload provided as a JSON string
  $ stackit beta telemetryrouter destination create --payload "{...}" --instance-id xxx

  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit beta telemetryrouter destination generate-payload --config-type s3 > ./payload.json
  <Modify payload in file, if needed>
  $ stackit beta telemetryrouter destination create --payload @./payload.json --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit beta telemetryrouter destination create"
      --instance-id string   ID of the TelemetryRouter instance
      --payload string       Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit beta telemetryrouter destination generate-payload")
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

