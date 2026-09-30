## stackit beta volume automation create

Creates a Volume Automation

### Synopsis

Creates a Volume Automation.
The input for the automation can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/automation-service/version/v1#tag/Volume-Automations/operation/CreateVolumeAutomation for information regarding the payload structure.


```
stackit beta volume automation create [flags]
```

### Examples

```
  Creates a Volume Automation with name "my-automation", template ID "xxx" and input from the file "./input-payload.json"
  $ stackit beta volume automation create --name my-automation --template-id xxx --input @./input-payload.json

  Creates a Volume Automation with name "my-automation", description "CLI Example", template ID "xxx" and input from the file "./input-payload.json"
  $ stackit beta volume automation create --name my-automation --description "CLI Example" --template-id xxx --input @./input-payload.json

  Creates a Volume Automation with name "my-automation", template ID "xxx", rrule schedule trigger "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR" and input from the file "./input-payload.json"
  $ stackit beta volume automation create --name my-automation --template-id xxx --trigger-schedule-rrule "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR" --input @./input-payload.json
```

### Options

```
      --description string              Description of the automation
  -h, --help                            Help for "stackit beta volume automation create"
      --input string                    Input for the automation (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).
      --name string                     Name of the automation
      --template-id string              Template ID which should be used for the automation
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

