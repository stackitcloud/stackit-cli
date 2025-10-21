## stackit beta intake runner create

Creates a new Intake Runner

### Synopsis

Creates a new Intake Runner.

```
stackit beta intake runner create [flags]
```

### Examples

```
  Create a new Intake Runner with a display name and message capacity limits
  $ stackit beta intake runner create --display-name my-runner --max-message-size-kib 1000 --max-messages-per-hour 5000

  Create a new Intake Runner with a description and labels
  $ stackit beta intake runner create --display-name my-runner --max-message-size-kib 1000 --max-messages-per-hour 5000 --description "Main runner for production" --labels="env=prod,team=billing"
```

### Options

```
      --description string          Description
      --display-name string         Display name
  -h, --help                        Help for "stackit beta intake runner create"
      --labels stringToString       Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2" (default [])
      --max-message-size-kib int    Maximum message size in KiB
      --max-messages-per-hour int   Maximum number of messages per hour
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta intake runner](./stackit_beta_intake_runner.md)	 - Provides functionality for Intake Runners

