## stackit ske cluster delete

Deletes a SKE cluster

### Synopsis

Deletes a STACKIT Kubernetes Engine (SKE) cluster.

```
stackit ske cluster delete CLUSTER_NAME [flags]
```

### Examples

```
  Delete a SKE cluster with name "my-cluster"
  $ stackit ske cluster delete my-cluster
```

### Options

```
  -h, --help   Help for "stackit ske cluster delete"
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

