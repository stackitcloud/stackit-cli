## stackit ske hibernate

Trigger hibernate for a SKE cluster

### Synopsis

Trigger hibernate for a STACKIT Kubernetes Engine (SKE) cluster.

```
stackit ske hibernate CLUSTER_NAME [flags]
```

### Examples

```
  Trigger hibernate for a SKE cluster with name "my-cluster"
  $ stackit ske cluster hibernate my-cluster
```

### Options

```
  -h, --help   Help for "stackit ske hibernate"
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

* [stackit ske](./stackit_ske.md)	 - Provides functionality for SKE

