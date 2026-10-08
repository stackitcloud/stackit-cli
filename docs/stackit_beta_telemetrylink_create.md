## stackit beta telemetrylink create

Creates a Telemetry Link with the specified resource

### Synopsis

Creates a Telemetry Link with the specified resource.

```
stackit beta telemetrylink create [flags]
```

### Examples

```
  Create a Telemetry Link with name "my-link" for project with ID "xxx" and for telemetry router with ID "yyy"
  $ stackit beta telemetrylink create --display-name my-link --resource-type project --resource-id xxx --telemetry-router-id yyy --access-token my-token

  Create a Telemetry Link with name "my-link" for folder with ID "xxx" and for telemetry router with ID "yyy"
  $ stackit beta telemetrylink create --display-name my-link --resource-type folder --resource-id xxx --telemetry-router-id yyy --access-token my-token

  Create a Telemetry Link with name "my-link" for organization with ID "xxx" and for telemetry router with ID "yyy"
  $ stackit beta telemetrylink create --display-name my-link --resource-type organization --resource-id xxx --telemetry-router-id yyy --access-token my-token

  Create a Telemetry Link with name "my-link" for organization with ID "xxx" and for telemetry router with ID "yyy", and disable routing
  $ stackit beta telemetrylink create --display-name my-link --resource-type organization --resource-id xxx --telemetry-router-id yyy --enabled=false --access-token my-token
```

### Options

```
      --access-token string          The access token of the Telemetry Router instance
      --description string           The description of the Telemetry Link resource
      --display-name string          The displayed name of the Telemetry Link resource
      --enabled                      Enable routing through the link to a telemetry router (default true)
  -h, --help                         Help for "stackit beta telemetrylink create"
      --resource-id string           STACKIT project ID, folder ID, or organization ID associated with the Telemetry Link resource
      --resource-type string         The resource type of the TelemetryLink resource, possible values are [organization folder project]
      --telemetry-router-id string   The ID of the telemetry-router to route the telemetry data
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

* [stackit beta telemetrylink](./stackit_beta_telemetrylink.md)	 - Provides functionality for Telemetry Link

