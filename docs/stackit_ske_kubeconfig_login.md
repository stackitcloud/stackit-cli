## stackit ske kubeconfig login

Login plugin for kubernetes clients

### Synopsis

Login plugin for kubernetes clients, that creates short-lived credentials to authenticate against a STACKIT Kubernetes Engine (SKE) cluster.
First you need to obtain a kubeconfig for use with the login command (first example).
Secondly you use the kubeconfig with your chosen Kubernetes client (second example), the client will automatically retrieve the credentials via the STACKIT CLI.

```
stackit ske kubeconfig login [flags]
```

### Examples

```
  Get a login kubeconfig for the SKE cluster with name "my-cluster". This kubeconfig does not contain any credentials and instead obtains valid credentials via the `stackit ske kubeconfig login` command.
  $ stackit ske kubeconfig create my-cluster --login

  Use the previously saved kubeconfig to authenticate to the SKE cluster, in this case with kubectl.
  $ kubectl cluster-info
  $ kubectl get pods
```

### Options

```
  -h, --help   Help for "stackit ske kubeconfig login"
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

* [stackit ske kubeconfig](./stackit_ske_kubeconfig.md)	 - Provides functionality for SKE kubeconfig

