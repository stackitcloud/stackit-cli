## stackit alpha sca application describe

Show details of a SCA application

### Synopsis

Show details of a STACKIT Container Applications (SCA) application.

```
stackit alpha sca application describe [flags]
```

### Examples

```
  Get details of a SCA application with ID "xxx" from an environment with ID "yyy"
  $ stackit alpha sca application describe xxx --environment-id yyy

  Get details of all SCA application with ID "xxx" from an environment with ID "yyy" in JSON format
  $ stackit alpha sca application describe xxx --environment-id yyy --output-format json
```

### Options

```
      --environment-id string   Environment ID (if not set uses default environment)
  -h, --help                    Help for "stackit alpha sca application describe"
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

* [stackit alpha sca application](./stackit_alpha_sca_application.md)	 - Provides functionality for SCA applications

