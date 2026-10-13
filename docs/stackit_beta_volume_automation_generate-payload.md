## stackit beta volume automation generate-payload

Generates a payload to update volume automation

### Synopsis

Generates a JSON payload with values to be used as --input for volume automation update.
See https://docs.api.stackit.cloud/documentation/automation-service/version/v1#tag/Volume-Automations/operation/PartialUpdateVolumeAutomation for information regarding the input structure.

```
stackit beta volume automation generate-payload [flags]
```

### Examples

```
  Generate a payload with values of a volume automation, and adapt it with custom values to different configuration options.
  $ stackit beta volume automation generate-payload --automation-id XXX --file-path ./input-payload.json
  <Modify payload in file>
  $ stackit beta volume automation update xxx --input @./input-payload.json

  Generate a payload with values of a volume automation, and preview it in the terminal.
  $ stackit beta volume automation generate-payload --automation-id XXX
```

### Options

```
      --automation-id string   Automation ID from which the input values should be used
  -f, --file-path string       If set, writes the payload to the given file. If unset, writes the payload to the standard output
  -h, --help                   Help for "stackit beta volume automation generate-payload"
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

* [stackit beta volume automation](./stackit_beta_volume_automation.md)	 - Provides functionality for Volume Automation

