## stackit beta telemetrylink describe

Shows details of a Telemetry Link within the specified resource

### Synopsis

Shows details of a Telemetry Link within the specified resource.

```
stackit beta telemetrylink describe [flags]
```

### Examples

```
  Get details of a Telemetry Link for project with ID "xxx"
  $ stackit beta telemetrylink describe --resource-type project --resource-id xxx

  Get details of a Telemetry Link for folder with ID "xxx"
  $ stackit beta telemetrylink describe --resource-type folder --resource-id xxx

  Get details of a Telemetry Link for organization with ID "xxx"
  $ stackit beta telemetrylink describe --resource-type organization --resource-id xxx

  Get details of a Telemetry Link for organization with ID "xxx" in JSON format
  $ stackit beta telemetrylink describe --resource-type organization --resource-id xxx --output-format json
```

### Options

```
  -h, --help                   Help for "stackit beta telemetrylink describe"
      --resource-id string     STACKIT project ID, folder ID, or organization ID associated with the Telemetry Link resource
      --resource-type string   The resource type of the TelemetryLink resource, possible values are [organization folder project]
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

