## stackit beta telemetryrouter access-token create

Creates a TelemetryRouter access token

### Synopsis

Creates a TelemetryRouter access token.

```
stackit beta telemetryrouter access-token create [flags]
```

### Examples

```
  Create an access token with the display name "access-token-1" for the TelemetryRouter instance "xxx"
  $ stackit beta telemetryrouter access-token create --display-name access-token-1 --instance-id xxx

  Create an access token with a description
  $ stackit beta telemetryrouter access-token create --display-name access-token-1 --instance-id xxx --description "Access token for service"

  Create an access token which expires in 30 days
  $ stackit beta telemetryrouter access-token create --display-name access-token-1 --instance-id xxx --ttl 30
```

### Options

```
      --description string    Description of the access token
      --display-name string   Display name for the access token
  -h, --help                  Help for "stackit beta telemetryrouter access-token create"
      --instance-id string    ID of the TelemetryRouter instance
      --ttl int32             Time-to-live (TTL) in days for the access token. If not set, the token will not expire
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

