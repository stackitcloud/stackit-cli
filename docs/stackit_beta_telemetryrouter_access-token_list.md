## stackit beta telemetryrouter access-token list

Lists all TelemetryRouter access tokens of an instance

### Synopsis

Lists all TelemetryRouter access tokens of an instance.

```
stackit beta telemetryrouter access-token list [flags]
```

### Examples

```
  Lists all access tokens of the TelemetryRouter instance "xxx"
  $ stackit beta telemetryrouter access-token list --instance-id xxx

  Lists all access tokens in JSON format
  $ stackit beta telemetryrouter access-token list --instance-id xxx --output-format json

  Lists up to 10 access tokens
  $ stackit beta telemetryrouter access-token list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit beta telemetryrouter access-token list"
      --instance-id string   ID of the TelemetryRouter instance
      --limit int            Maximum number of entries to list
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

* [stackit beta telemetryrouter access-token](./stackit_beta_telemetryrouter_access-token.md)	 - Provides functionality for TelemetryRouter access tokens

