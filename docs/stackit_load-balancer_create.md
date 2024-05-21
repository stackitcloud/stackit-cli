## stackit load-balancer create

Creates a Load Balancer

### Synopsis

Creates a Load Balancer.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/load-balancer/version/v1#tag/Load-Balancer/operation/APIService_CreateLoadBalancer for information regarding the payload structure.

```
stackit load-balancer create [flags]
```

### Examples

```
  Create a load balancer using an API payload sourced from the file "./payload.json"
  $ stackit load-balancer create --payload @./payload.json

  Create a load balancer using an API payload provided as a JSON string
  $ stackit load-balancer create --payload "{...}"

  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit load-balancer generate-payload > ./payload.json
  <Modify payload in file>
  $ stackit load-balancer create --payload @./payload.json
```

### Options

```
  -h, --help             Help for "stackit load-balancer create"
      --payload string   Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).
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

