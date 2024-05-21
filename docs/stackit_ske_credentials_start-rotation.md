## stackit ske credentials start-rotation

Starts the rotation of the credentials associated to a SKE cluster

### Synopsis

Starts the rotation of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster.

This is step 1 of a 2-step process to rotate all SKE cluster credentials. Tasks accomplished in this phase include:
  - Rolling recreation of all worker nodes
  - A new Certificate Authority (CA) will be established and incorporated into the existing CA bundle.
  - A new etcd encryption key is generated and added to the Certificate Authority (CA) bundle.
  - A new signing key will be generated for the service account and added to the Certificate Authority (CA) bundle.
  - The kube-apiserver will rewrite all secrets in the cluster, encrypting them with the new encryption key.
The old CA, encryption key and signing key will be retained until the rotation is completed.

After completing the rotation of credentials, you can generate a new kubeconfig file by running:
  $ stackit ske kubeconfig create my-cluster
Complete the rotation by running:
  $ stackit ske credentials complete-rotation my-cluster
For more information, visit: https://docs.stackit.cloud/stackit/en/how-to-rotate-ske-credentials-200016334.html

```
stackit ske credentials start-rotation CLUSTER_NAME [flags]
```

### Examples

```
  Start the rotation of the credentials associated to the SKE cluster with name "my-cluster"
  $ stackit ske credentials start-rotation my-cluster

  Flow of the 2-step process to rotate all SKE cluster credentials, including generating a new kubeconfig file
  $ stackit ske credentials start-rotation my-cluster
  $ stackit ske kubeconfig create my-cluster
  $ stackit ske credentials complete-rotation my-cluster
```

### Options

```
  -h, --help   Help for "stackit ske credentials start-rotation"
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

* [stackit ske credentials](./stackit_ske_credentials.md)	 - Provides functionality for SKE credentials

