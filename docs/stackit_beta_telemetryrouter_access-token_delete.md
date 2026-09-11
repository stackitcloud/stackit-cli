## stackit beta telemetryrouter access-token delete

Deletes a TelemetryRouter access token

### Synopsis

Deletes a TelemetryRouter access token.

```
stackit beta telemetryrouter access-token delete ACCESS_TOKEN_ID [flags]
```

### Examples

```
  Delete access token with ID "xxx" for the TelemetryRouter instance "yyy"
  $ stackit beta telemetryrouter access-token delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit beta telemetryrouter access-token delete"
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

* [stackit beta telemetryrouter access-token](./stackit_beta_telemetryrouter_access-token.md)	 - Provides functionality for TelemetryRouter access tokens

