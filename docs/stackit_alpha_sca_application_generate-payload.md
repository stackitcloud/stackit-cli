## stackit alpha sca application generate-payload

Generates a payload to create/update SCA applications

### Synopsis

Generates a JSON payload with values to be used as --payload input for application creation or update.
See https://docs.api.stackit.cloud/documentation/sca/version/v1alpha#tag/Applications/operation/Applications_CreateApplication for information regarding the payload structure.

```
stackit alpha sca application generate-payload [flags]
```

### Examples

```
  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit alpha sca application generate-payload --file-path ./payload.json
  <Modify payload in file, if needed>
  $ stackit alpha sca application create-from-payload --name application-name --payload @./payload.json

  Generate a payload with values of an application, and adapt it with custom values for the different configuration options
  $ stackit alpha sca application generate-payload --application-id xxx --file-path ./payload.json
  <Modify payload in file>
  $ stackit alpha sca application update-from-payload --payload @./payload.json

  Generate a payload with values of an application, and preview it in the terminal
  $ stackit alpha sca application generate-payload --application-id xxx
```

### Options

```
      --application-id string   
      --environment-id string   
  -f, --file-path string        If set, writes the payload to the given file. If unset, writes the payload to the standard output
  -h, --help                    Help for "stackit alpha sca application generate-payload"
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

