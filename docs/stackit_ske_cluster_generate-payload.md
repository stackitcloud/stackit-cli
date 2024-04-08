## stackit ske cluster generate-payload

Generates a payload to create/update SKE clusters

### Synopsis

Generates a JSON payload with values to be used as --payload input for cluster creation or update.
See https://docs.api.stackit.cloud/documentation/ske/version/v1#tag/Cluster/operation/SkeService_CreateOrUpdateCluster for information regarding the payload structure.

```
stackit ske cluster generate-payload [flags]
```

### Examples

```
  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit ske cluster generate-payload > ./payload.json
  <Modify payload in file, if needed>
  $ stackit ske cluster create my-cluster --payload @./payload.json

  Generate a payload with values of a cluster, and adapt it with custom values for the different configuration options
  $ stackit ske cluster generate-payload --cluster-name my-cluster > ./payload.json
  <Modify payload in file>
  $ stackit ske cluster update my-cluster --payload @./payload.json
```

### Options

```
  -n, --cluster-name string   If set, generates the payload with the current state of the given cluster. If unset, generates the payload with default values
  -h, --help                  Help for "stackit ske cluster generate-payload"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit ske cluster](./stackit_ske_cluster.md)	 - Provides functionality for SKE cluster

