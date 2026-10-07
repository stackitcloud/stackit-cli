## stackit alpha sca application create-from-payload

Create a SCA application from payload

### Synopsis

Create a STACKIT Container Applications (SCA) application from payload.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/sca/version/v1alpha#tag/Applications/operation/Applications_CreateApplication for information regarding the payload structure.

```
stackit alpha sca application create-from-payload [flags]
```

### Examples

```
  Create a SCA application using an API payload sourced from the file "./payload.json"
  $ stackit alpha sca cluster create-from-payload --name application-name --payload @./payload.json

  Create a SCA application using an API payload provided as a JSON string
  $ stackit alpha sca cluster create-from-payload --name application-name --payload "{...}"

  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit alpha sca application generate-payload --file-path ./payload.json
  <Modify payload in file, if needed>
  $ stackit alpha sca application create-from-payload --name application-name --payload @./payload.json
```

### Options

```
      --environment-id string   Environment ID (uses default environment if not set)
  -h, --help                    Help for "stackit alpha sca application create-from-payload"
      --name string             Application display name
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

