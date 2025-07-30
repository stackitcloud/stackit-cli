## stackit ske cluster update

Updates a SKE cluster

### Synopsis

Updates a STACKIT Kubernetes Engine (SKE) cluster.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/ske/version/v1#tag/Cluster/operation/SkeService_CreateOrUpdateCluster for information regarding the payload structure.

```
stackit ske cluster update CLUSTER_NAME [flags]
```

### Examples

```
  Update a SKE cluster using an API payload sourced from the file "./payload.json"
  $ stackit ske cluster update my-cluster --payload @./payload.json

  Update a SKE cluster using an API payload provided as a JSON string
  $ stackit ske cluster update my-cluster --payload "{...}"

  Generate a payload with the current values of a cluster, and adapt it with custom values for the different configuration options
  $ stackit ske cluster generate-payload --cluster-name my-cluster > ./payload.json
  <Modify payload in file>
  $ stackit ske cluster update my-cluster --payload @./payload.json
```

### Options

```
  -h, --help             Help for "stackit ske cluster update"
      --payload string   Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json
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

* [stackit ske cluster](./stackit_ske_cluster.md)	 - Provides functionality for SKE cluster

