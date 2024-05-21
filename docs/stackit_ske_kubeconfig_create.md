## stackit ske kubeconfig create

Creates a kubeconfig for an SKE cluster

### Synopsis

Creates a kubeconfig for a STACKIT Kubernetes Engine (SKE) cluster.

By default the kubeconfig is created in the .kube folder, in the user's home directory. The kubeconfig file will be overwritten if it already exists.
You can override this behavior by specifying a custom filepath with the --filepath flag.
An expiration time can be set for the kubeconfig. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is 1h.
Note that the format is <value><unit>, e.g. 30d for 30 days and you can't combine units.

```
stackit ske kubeconfig create CLUSTER_NAME [flags]
```

### Examples

```
  Create a kubeconfig for the SKE cluster with name "my-cluster"
  $ stackit ske kubeconfig create my-cluster

  Get a login kubeconfig for the SKE cluster with name "my-cluster". This kubeconfig does not contain any credentials and instead obtains valid credentials via the `stackit ske kubeconfig login` command.
  $ stackit ske kubeconfig create my-cluster --login

  Create a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 30 days
  $ stackit ske kubeconfig create my-cluster --expiration 30d

  Create a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 2 months
  $ stackit ske kubeconfig create my-cluster --expiration 2M

  Create a kubeconfig for the SKE cluster with name "my-cluster" in a custom filepath
  $ stackit ske kubeconfig create my-cluster --filepath /path/to/config
```

### Options

```
  -e, --expiration string   Expiration time for the kubeconfig in seconds(s), minutes(m), hours(h), days(d) or months(M). Example: 30d. By default, expiration time is 1h
      --filepath string     Path to create the kubeconfig file. By default, the kubeconfig is created as 'config' in the .kube folder, in the user's home directory.
  -h, --help                Help for "stackit ske kubeconfig create"
  -l, --login               Create a login kubeconfig that obtains valid credentials via the STACKIT CLI. This flag is mutually exclusive with the expiration flag.
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

* [stackit ske kubeconfig](./stackit_ske_kubeconfig.md)	 - Provides functionality for SKE kubeconfig

