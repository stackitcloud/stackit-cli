## stackit alpha sca application update-from-payload

Update a SCA application from payload

### Synopsis

Update a STACKIT Container Applications (SCA) application from payload.

```
stackit alpha sca application update-from-payload [flags]
```

### Examples

```
  Update a SCA application using an API payload sourced from the file "./payload.json"
  $ stackit alpha sca application update-from-payload my-application-id --payload @./payload.json

  Update a SCA application using an API payload provided as a JSON string
  $ stackit alpha sca application update-from-payload my-application-id --payload "{...}"

  Generate a payload with the current values of an application, and adapt it with custom values for the different configuration options
  $ stackit alpha sca application generate-payload --application-id application-id > ./payload.json
  <Modify payload in file>
  $ stackit alpha sca application update-from-payload application-id --payload @./payload.json
```

### Options

```
      --environment-id string   Environment ID (uses default environment if not set)
  -h, --help                    Help for "stackit alpha sca application update-from-payload"
      --payload string          Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).
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

