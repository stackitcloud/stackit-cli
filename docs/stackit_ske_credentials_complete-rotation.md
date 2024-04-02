## stackit ske credentials complete-rotation

Completes the rotation of the credentials associated to a SKE cluster

### Synopsis

Completes the rotation of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster.

This is step 2 of a 2-step process to rotate all SKE cluster credentials. Tasks accomplished in this phase include:
  - The old certification authority will be dropped from the package.
  - The old signing key for the service account will be dropped from the bundle.
To ensure continued access to the Kubernetes cluster, please update your kubeconfig with the new credentials:
  $ stackit ske kubeconfig create my-cluster

If you haven't, please start the process by running:
  $ stackit ske credentials start-rotation my-cluster

```
stackit ske credentials complete-rotation CLUSTER_NAME [flags]
```

### Examples

```
  Complete the rotation of the credentials associated to the SKE cluster with name "my-cluster"
  $ stackit ske credentials complete-rotation my-cluster

  Flow of the 2-step process to rotate all SKE cluster credentials, including generating a new kubeconfig file
  $ stackit ske credentials start-rotation my-cluster
  $ stackit ske kubeconfig create my-cluster
  $ stackit ske credentials complete-rotation my-cluster
```

### Options

```
  -h, --help   Help for "stackit ske credentials complete-rotation"
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

