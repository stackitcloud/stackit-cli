## stackit beta telemetryrouter access-token update

Updates a TelemetryRouter access token

### Synopsis

Updates a TelemetryRouter access token.

```
stackit beta telemetryrouter access-token update ACCESS_TOKEN_ID [flags]
```

### Examples

```
  Update the display name of the access token with ID "xxx"
  $ stackit beta telemetryrouter access-token update xxx --instance-id yyy --display-name access-token-1

  Update the description of the access token with ID "xxx"
  $ stackit beta telemetryrouter access-token update xxx --instance-id yyy --description "Access token for service"
```

### Options

```
      --description string    Description of the access token
      --display-name string   Display name for the access token
  -h, --help                  Help for "stackit beta telemetryrouter access-token update"
      --instance-id string    ID of the TelemetryRouter instance
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

