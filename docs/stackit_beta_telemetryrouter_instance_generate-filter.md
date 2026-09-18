## stackit beta telemetryrouter instance generate-filter

Generates a filter configuration to use as --filter input for TelemetryRouter instances

### Synopsis

Generates a JSON filter configuration to be used as "--filter" input for "instance create" or "instance update".
If "--instance-id" is set, the current filter configuration of that instance is returned, to be adapted with custom values.
If unset, a default example filter configuration is returned.

```
stackit beta telemetryrouter instance generate-filter [flags]
```

### Examples

```
  Generate a default example filter configuration, and adapt it with custom values
  $ stackit beta telemetryrouter instance generate-filter --file-path ./filter.json
  <Modify filter in file, if needed>
  $ stackit beta telemetryrouter instance create --display-name "my-instance" --filter @./filter.json

  Generate a filter configuration with the current values of TelemetryRouter instance "xxx", and adapt it with custom values
  $ stackit beta telemetryrouter instance generate-filter --instance-id xxx --file-path ./filter.json
  <Modify filter in file>
  $ stackit beta telemetryrouter instance update xxx --filter @./filter.json

  Generate a filter configuration with the current values of TelemetryRouter instance "xxx", and preview it in the terminal
  $ stackit beta telemetryrouter instance generate-filter --instance-id xxx
```

### Options

```
  -f, --file-path string     If set, writes the filter configuration to the given file. If unset, writes it to the standard output
  -h, --help                 Help for "stackit beta telemetryrouter instance generate-filter"
      --instance-id string   If set, generates a filter configuration with the current values of the given TelemetryRouter instance. If unset, generates a default example filter configuration
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

* [stackit beta telemetryrouter instance](./stackit_beta_telemetryrouter_instance.md)	 - Provides functionality for TelemetryRouter instances

