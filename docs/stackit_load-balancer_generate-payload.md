## stackit load-balancer generate-payload

Generates a payload to create/update a Load Balancer

### Synopsis

Generates a JSON payload with values to be used as --payload input for load balancer creation or update.
See https://docs.api.stackit.cloud/documentation/load-balancer/version/v1#tag/Load-Balancer/operation/APIService_CreateLoadBalancer for information regarding the payload structure.

```
stackit load-balancer generate-payload [flags]
```

### Examples

```
  Generate a payload, and adapt it with custom values for the different configuration options
  $ stackit load-balancer generate-payload --file-path ./payload.json
  <Modify payload in file, if needed>
  $ stackit load-balancer create --payload @./payload.json

  Generate a payload with values of an existing load balancer, and adapt it with custom values for the different configuration options
  $ stackit load-balancer generate-payload --lb-name xxx --file-path ./payload.json
  <Modify payload in file>
  $ stackit load-balancer update xxx --payload @./payload.json

  Generate a payload with values of an existing load balancer, and preview it in the terminal
  $ stackit load-balancer generate-payload --lb-name xxx
```

### Options

```
  -f, --file-path string   If set, writes the payload to the given file. If unset, writes the payload to the standard output
  -h, --help               Help for "stackit load-balancer generate-payload"
  -n, --lb-name string     If set, generates the payload with the current values of the given load balancer. If unset, generates the payload with empty values
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

* [stackit load-balancer](./stackit_load-balancer.md)	 - Provides functionality for Load Balancer

