## stackit ske cluster create

Creates an SKE cluster

### Synopsis

Creates a STACKIT Kubernetes Engine (SKE) cluster.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/ske/version/v1#tag/Cluster/operation/SkeService_CreateOrUpdateCluster for information regarding the payload structure.

```
stackit ske cluster create CLUSTER_NAME [flags]
```

### Examples

```
  Create an SKE cluster using default configuration
  $ stackit ske cluster create my-cluster

  Create an SKE cluster using an API payload sourced from the file "./payload.json"
  $ stackit ske cluster create my-cluster --payload @./payload.json

  Create an SKE cluster using an API payload provided as a JSON string
  $ stackit ske cluster create my-cluster --payload "{...}"

  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit ske cluster generate-payload > ./payload.json
  <Modify payload in file, if needed>
  $ stackit ske cluster create my-cluster --payload @./payload.json
```

### Options

```
  -h, --help             Help for "stackit ske cluster create"
      --payload string   Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit ske cluster generate-payload")
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

* [stackit ske cluster](./stackit_ske_cluster.md)	 - Provides functionality for SKE cluster

