## stackit ske cluster describe

Shows details of a SKE cluster

### Synopsis

Shows details of a STACKIT Kubernetes Engine (SKE) cluster.

```
stackit ske cluster describe CLUSTER_NAME [flags]
```

### Examples

```
  Get details of a SKE cluster with name "my-cluster"
  $ stackit ske cluster describe my-cluster

  Get details of a SKE cluster with name "my-cluster" in JSON format
  $ stackit ske cluster describe my-cluster --output-format json
```

### Options

```
  -h, --help   Help for "stackit ske cluster describe"
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

