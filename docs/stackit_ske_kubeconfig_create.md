## stackit ske kubeconfig create

Creates or update a kubeconfig for an SKE cluster

### Synopsis

Creates a kubeconfig for a STACKIT Kubernetes Engine (SKE) cluster, if the config exits in the kubeconfig file the information will be updated.

By default, the kubeconfig information of the SKE cluster is merged into the default kubeconfig file of the current user. If the kubeconfig file doesn't exist, a new one will be created.
You can override this behavior by specifying a custom filepath with the --filepath flag.

An expiration time can be set for the kubeconfig. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is 1h.

Note that the format is <value><unit>, e.g. 30d for 30 days and you can't combine units.

```
stackit ske kubeconfig create CLUSTER_NAME [flags]
```

### Examples

```
  Create a kubeconfig for the SKE cluster with name "my-cluster. If the config exits in the kubeconfig file the information will be updated."
  $ stackit ske kubeconfig create my-cluster

  Get a login kubeconfig for the SKE cluster with name "my-cluster". This kubeconfig does not contain any credentials and instead obtains valid credentials via the `stackit ske kubeconfig login` command.
  $ stackit ske kubeconfig create my-cluster --login

  Create o kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 30 days. If the config exits in the kubeconfig file the information will be updated.
  $ stackit ske kubeconfig create my-cluster --expiration 30d

  Create or update a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 2 months. If the config exits in the kubeconfig file the information will be updated.
  $ stackit ske kubeconfig create my-cluster --expiration 2M

  Create or update a kubeconfig for the SKE cluster with name "my-cluster" in a custom filepath. If the config exits in the kubeconfig file the information will be updated.
  $ stackit ske kubeconfig create my-cluster --filepath /path/to/config

  Get a kubeconfig for the SKE cluster with name "my-cluster" without writing it to a file and format the output as json
  $ stackit ske kubeconfig create my-cluster --disable-writing --output-format json
```

### Options

```
      --disable-writing     Disable the writing of kubeconfig. Set the output format to json or yaml using the --output-format flag to display the kubeconfig.
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
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit ske kubeconfig](./stackit_ske_kubeconfig.md)	 - Provides functionality for SKE kubeconfig

