## stackit intake runner update

Updates an Intake Runner

### Synopsis

Updates an Intake Runner. Only the specified fields are updated.

```
stackit intake runner update RUNNER_ID [flags]
```

### Examples

```
  Update the display name of an Intake Runner with ID "xxx"
  $ stackit intake runner update xxx --display-name "new-runner-name"

  Update the message capacity limits for an Intake Runner with ID "xxx"
  $ stackit intake runner update xxx --max-message-size-kib 2000 --max-messages-per-hour 10000

  Clear the labels of an Intake Runner with ID "xxx" by providing an empty value
  $ stackit intake runner update xxx --labels ""
```

### Options

```
      --description string          Description
      --display-name string         Display name
  -h, --help                        Help for "stackit intake runner update"
      --labels string               Labels in key=value format. To clear all labels, provide an empty string, e.g. --labels ""
      --max-message-size-kib int    Maximum message size in KiB. Note: Overall message capacity cannot be decreased.
      --max-messages-per-hour int   Maximum number of messages per hour. Note: Overall message capacity cannot be decreased.
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

* [stackit intake runner](./stackit_intake_runner.md)	 - Provides functionality for Intake Runners

