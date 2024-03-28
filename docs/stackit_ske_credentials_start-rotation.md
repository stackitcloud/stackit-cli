## stackit ske credentials start-rotation

Starts the rotation of the credentials associated to a SKE cluster

### Synopsis

Starts the rotation of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster. This is step 1 of a two-step process. 
Complete the rotation using the 'stackit ske credentials complete-rotation' command.

```
stackit ske credentials start-rotation CLUSTER_NAME [flags]
```

### Examples

```
  Start the rotation of the credentials associated to the SKE cluster with name "my-cluster"
  $ stackit ske credentials start-rotation my-cluster
```

### Options

```
  -h, --help   Help for "stackit ske credentials start-rotation"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit ske credentials](./stackit_ske_credentials.md)	 - Provides functionality for SKE credentials

