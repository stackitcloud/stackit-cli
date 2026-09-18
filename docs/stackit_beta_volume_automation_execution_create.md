## stackit beta volume automation execution create

Creates a new Volume Automation Execution

### Synopsis

Creates a new Volume Automation Execution.

```
stackit beta volume automation execution create [flags]
```

### Examples

```
  Create a new Volume Automation Execution for Automation with ID "xxx"
  $ stackit beta volume automation execution create --automation-id xxx
```

### Options

```
      --automation-id string   ID of the automation which should be executed
  -h, --help                   Help for "stackit beta volume automation execution create"
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

* [stackit beta volume automation execution](./stackit_beta_volume_automation_execution.md)	 - Provides functionality for Volume Automation Execution

