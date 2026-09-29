## stackit beta volume automation list

List all Volume Automations

### Synopsis

List all Volume Automations.

```
stackit beta volume automation list [flags]
```

### Examples

```
  List all Volume Automations
  $ stackit beta volume automation list

  List up to 10 Volume Automations
  $ stackit beta volume automation list --limit 10
```

### Options

```
  -h, --help        Help for "stackit beta volume automation list"
      --limit int   Maximum number of entries to list
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

