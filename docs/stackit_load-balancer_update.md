## stackit load-balancer update

Updates a Load Balancer

### Synopsis

Updates a load balancer.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/load-balancer/version/v1#tag/Load-Balancer/operation/APIService_UpdateLoadBalancer for information regarding the payload structure.

```
stackit load-balancer update LOAD_BALANCER_NAME [flags]
```

### Examples

```
  Update a load balancer with name "xxx", using an API payload sourced from the file "./payload.json"
  $ stackit load-balancer update xxx --payload @./payload.json

  Update a load balancer with name "xxx", using an API payload provided as a JSON string
  $ stackit load-balancer update xxx --payload "{...}"

  Generate a payload with the current values of an existing load balancer xxx, and adapt it with custom values for the different configuration options
  $ stackit load-balancer generate-payload --lb-name xxx > ./payload.json
  <Modify payload in file>
  $ stackit load-balancer update xxx --payload @./payload.json
```

### Options

```
  -h, --help             Help for "stackit load-balancer update"
      --payload string   Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit load-balancer](./stackit_load-balancer.md)	 - Provides functionality for Load Balancer

