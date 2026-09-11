## stackit beta telemetryrouter destination update

Updates a Telemetry Router destination

### Synopsis

Updates a destination of a Telemetry Router instance.
Only the fields provided in the payload are updated, omitted fields are left unchanged.
The payload can be provided as a JSON string or a file path prefixed with "@".
If the "config" field is set, at most one of "openTelemetry" or "s3" may be set, matching the "configType".

```
stackit beta telemetryrouter destination update DESTINATION_ID [flags]
```

### Examples

```
  Update a destination with ID "xxx" from Telemetry Router instance "yyy", using an API payload sourced from the file "./payload.json"
  $ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload @./payload.json

  Update a destination with ID "xxx" from Telemetry Router instance "yyy", using an API payload provided as a JSON string
  $ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload "{...}"

  Generate a payload with the current values of a destination, and adapt it with custom values for the different configuration options
  $ stackit beta telemetryrouter destination generate-payload --destination-id xxx --instance-id yyy --file-path ./payload.json
  <Modify payload in file>
  $ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload @./payload.json
```

### Options

```
  -h, --help[=true|false]    Help for "stackit beta telemetryrouter destination update"
      --instance-id string   ID of the Telemetry Router instance
      --payload string       Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit beta telemetryrouter destination](./stackit_beta_telemetryrouter_destination.md)	 - Provides functionality for Telemetry Router destinations

