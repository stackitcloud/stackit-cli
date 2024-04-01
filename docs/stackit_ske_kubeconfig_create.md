## stackit ske kubeconfig create

Creates a kubeconfig for an SKE cluster

### Synopsis

Creates a kubeconfig for a STACKIT Kubernetes Engine (SKE) cluster.
By default the kubeconfig is created in the .kube folder, in the user's home directory. The kubeconfig file will be overwritten if it already exists.

```
stackit ske kubeconfig create CLUSTER_NAME [flags]
```

### Examples

```
  Create a kubeconfig for the SKE cluster with name "my-cluster"
  $ stackit ske kubeconfig create my-cluster

  Create a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 30 days
  $ stackit ske kubeconfig create my-cluster --expiration 30d

  Create a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 2 months
  $ stackit ske kubeconfig create my-cluster --expiration 2M

  Create a kubeconfig for the SKE cluster with name "my-cluster" in a custom location
  $ stackit ske kubeconfig create my-cluster --location /path/to/config
```

### Options

```
  -e, --expiration string   Expiration time for the kubeconfig in seconds(s), minutes(m), hours(h), days(d) or months(M). Example: 30d. By default, expiration time is 1h
  -h, --help                Help for "stackit ske kubeconfig create"
      --location string     Folder location to store the kubeconfig file. By default, the kubeconfig is created in the .kube folder, in the user's home directory.
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit ske kubeconfig](./stackit_ske_kubeconfig.md)	 - Provides functionality for SKE kubeconfig

