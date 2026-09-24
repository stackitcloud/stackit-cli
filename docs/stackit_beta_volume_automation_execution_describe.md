## stackit beta volume automation execution describe

Shows details of a Volume Automation Execution

### Synopsis

Shows details of a Volume Automation Execution.

```
stackit beta volume automation execution describe EXECUTION_ID [flags]
```

### Examples

```
  Describe the Volume Automation Execution with ID "xxx" for Automation with ID "yyy"
  $ stackit beta volume automation execution describe xxx --automation-id yyy
```

### Options

```
      --automation-id string   Automation ID
  -h, --help                   Help for "stackit beta volume automation execution describe"
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

