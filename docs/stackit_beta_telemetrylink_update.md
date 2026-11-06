## stackit beta telemetrylink update

Updates a Telemetry Link within the specified resource

### Synopsis

Updates a Telemetry Link within the specified resource.

```
stackit beta telemetrylink update [flags]
```

### Examples

```
  Update a Telemetry Link with name "my-new-link" for project with ID "xxx"
  $ stackit beta telemetrylink update --display-name my-new-link --resource-type project --resource-id xxx

  Update a Telemetry Link with new access token for folder with ID "xxx"
  $ stackit beta telemetrylink update --resource-type folder --resource-id xxx --access-token my-new-token

  Update a Telemetry Link with new telemetry router ID "xxx" for organization with ID "yyy"
  $ stackit beta telemetrylink update --resource-type organization --resource-id yyy --telemetry-router-id xxx
```

### Options

```
      --access-token string          The access token of the Telemetry Router instance
      --description string           The description of the Telemetry Link resource
      --display-name string          The displayed name of the Telemetry Link resource
      --enabled                      Enable routing through the link to a telemetry router (default true)
  -h, --help                         Help for "stackit beta telemetrylink update"
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

