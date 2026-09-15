## stackit beta telemetryrouter destination update

Updates a TelemetryRouter destination

### Synopsis

Updates a destination of a TelemetryRouter instance.
The payload can be provided as a JSON string or a file path prefixed with "@".
The "config" field is a discriminated union: exactly one of "openTelemetry" or "s3" must be set, matching the "configType".

```
stackit beta telemetryrouter destination update DESTINATION_ID [flags]
```

### Examples

```
  Update a destination with ID "xxx" from TelemetryRouter instance "yyy", using an API payload sourced from the file "./payload.json"
  $ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload @./payload.json

  Update a destination with ID "xxx" from TelemetryRouter instance "yyy", using an API payload provided as a JSON string
  $ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload "{...}"

  Generate a payload with the current values of a destination, and adapt it with custom values for the different configuration options
  $ stackit beta telemetryrouter destination generate-payload --destination-id xxx --instance-id yyy > ./payload.json
  <Modify payload in file>
  $ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload @./payload.json
```

### Options

```
  -h, --help                 Help for "stackit beta telemetryrouter destination update"
      --instance-id string   ID of the TelemetryRouter instance
      --payload string       Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json
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

