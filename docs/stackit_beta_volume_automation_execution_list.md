## stackit beta volume automation execution list

List all Volume Automation Executions

### Synopsis

List all Volume Automation Executions.

```
stackit beta volume automation execution list [flags]
```

### Examples

```
  List all Volume Automation Executions for Automation with ID "xxx"
  $ stackit beta volume automation execution list --automation-id xxx

  List up to 10 Volume Automation Executions for Automation with ID "xxx"
  $ stackit beta volume automation execution list --automation-id xxx --limit 10
```

### Options

```
      --automation-id string   Automation ID
  -h, --help                   Help for "stackit beta volume automation execution list"
      --limit int              Maximum number of entries to list
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

