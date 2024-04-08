## stackit ske credentials rotate

Rotates credentials associated to a SKE cluster

### Synopsis

Rotates credentials associated to a STACKIT Kubernetes Engine (SKE) cluster. The old credentials will be invalid after the operation.

```
stackit ske credentials rotate CLUSTER_NAME [flags]
```

### Examples

```
  Rotate credentials associated to the SKE cluster with name "my-cluster"
  $ stackit ske credentials rotate my-cluster
```

### Options

```
  -h, --help   Help for "stackit ske credentials rotate"
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

* [stackit ske credentials](./stackit_ske_credentials.md)	 - Provides functionality for SKE credentials

