## stackit ske wakeup

Trigger wakeup from hibernation for a SKE cluster

### Synopsis

Trigger wakeup from hibernation for a STACKIT Kubernetes Engine (SKE) cluster.

```
stackit ske wakeup CLUSTER_NAME [flags]
```

### Examples

```
  Trigger wakeup from hibernation for a SKE cluster with name "my-cluster"
  $ stackit ske cluster wakeup my-cluster
```

### Options

```
  -h, --help   Help for "stackit ske wakeup"
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

