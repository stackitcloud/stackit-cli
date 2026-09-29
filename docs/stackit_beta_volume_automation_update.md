## stackit beta volume automation update

Updates a Volume Automation

### Synopsis

Updates a Volume Automation.
The input for the automation can be provided as a JSON string or a file path prefixed with "@".
Updates are always applied as a patch. Omitted fields in the JSON will be ignored and remain unchanged after the update. Fields can be removed by setting them explicitly to null.
See https://docs.api.stackit.cloud/documentation/automation-service/version/v1#tag/Volume-Automations/operation/PartialUpdateVolumeAutomation for information regarding the payload structure.


```
stackit beta volume automation update AUTOMATION_ID [flags]
```

### Examples

```
  Update volume automation with ID "xxx" with input from the file "./input-payload.json"
  $ stackit beta volume automation update xxx --input @./input-payload.json

  Update volume automation with ID "xxx" with name "my-updated-automation"
  $ stackit beta volume automation update xxx --name "my-updated-automation"

  Update volume automation with ID "xxx" with description "Updated automation from STACKIT CLI"
  $ stackit beta volume automation update xxx --description "Updated automation from STACKIT CLI"

  Update volume automation with ID "xxx" with disabling trigger schedule rrule
  $ stackit beta volume automation update xxx --disable-trigger-schedule

  Update volume automation with ID "xxx" with trigger schedule rrule "FREQ=WEEKLY;BYDAY=MO"
  $ stackit beta volume automation update xxx --trigger-schedule-rrule "FREQ=WEEKLY;BYDAY=MO"
```

### Options

```
      --description string              Description of the automation
      --disable-trigger-schedule        If set, trigger schedule will be disabled
  -h, --help                            Help for "stackit beta volume automation update"
      --input string                    Input for the automation (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).
      --name string                     Name of the automation
      --trigger-schedule-rrule string   Trigger schedule (RRULE) for the automation
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

