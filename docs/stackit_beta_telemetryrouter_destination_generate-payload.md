## stackit beta telemetryrouter destination generate-payload

Generates a payload to create/update TelemetryRouter destinations

### Synopsis

Generates a JSON payload with values to be used as --payload input for destination creation or update.
This command can be used to generate a payload to update an existing destination or to create a new destination.
To update an existing destination, provide the destination ID and the instance ID of the TelemetryRouter instance.
To obtain a default payload to create a new destination, run the command with the "--config-type" flag set to either "opentelemetry" (default) or "s3".
Note that the default values provided, such as the URI, bucket name or credentials, should be adapted to your use case.

```
stackit beta telemetryrouter destination generate-payload [flags]
```

### Examples

```
  Generate a create payload with default values for an OpenTelemetry destination, and adapt it with custom values
  $ stackit beta telemetryrouter destination generate-payload --file-path ./payload.json
  <Modify payload in file, if needed>
  $ stackit beta telemetryrouter destination create --instance-id xxx --payload @./payload.json

  Generate a create payload with default values for an S3 destination
  $ stackit beta telemetryrouter destination generate-payload --config-type s3 --file-path ./payload.json

  Generate an update payload with the values of an existing destination "yyy" for TelemetryRouter instance "xxx", and adapt it with custom values
  $ stackit beta telemetryrouter destination generate-payload --destination-id yyy --instance-id xxx --file-path ./payload.json
  <Modify payload in file>
  $ stackit beta telemetryrouter destination update yyy --instance-id xxx --payload @./payload.json

  Generate an update payload with the values of an existing destination "yyy" for TelemetryRouter instance "xxx", and preview it in the terminal
  $ stackit beta telemetryrouter destination generate-payload --destination-id yyy --instance-id xxx
```

### Options

```
      --config-type string      Type of the default create payload to generate, one of "opentelemetry" or "s3". Only relevant if "--destination-id" is unset (default "opentelemetry")
      --destination-id string   If set, generates an update payload with the current state of the given destination. If unset, generates a create payload with default values
  -f, --file-path string        If set, writes the payload to the given file. If unset, writes the payload to the standard output
  -h, --help                    Help for "stackit beta telemetryrouter destination generate-payload"
      --instance-id string      ID of the TelemetryRouter instance
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

