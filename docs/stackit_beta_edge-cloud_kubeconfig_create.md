## stackit beta edge-cloud kubeconfig create

Creates or updates a local kubeconfig file of an Edge Cloud instance

### Synopsis

Creates or updates a local kubeconfig file of a STACKIT Edge Cloud (STEC) instance. If the config exists in the kubeconfig file, the information will be updated.

By default, the kubeconfig information of the edge instance is merged into the current kubeconfig file which is determined by Kubernetes client logic. If the kubeconfig file doesn't exist, a new one will be created.
You can override this behavior by specifying a custom filepath with the --filepath flag or disable writing with the --disable-writing flag.
An expiration time can be set for the kubeconfig. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is 3600 seconds.
Note: the format for the duration is <value><unit>, e.g. 30d for 30 days. You may not combine units.

```
stackit beta edge-cloud kubeconfig create [flags]
```

### Examples

```
  Create or update a kubeconfig for the Edge Cloud instance with instance ID "xxx". If the config exists in the kubeconfig file, the information will be updated.
  $ stackit beta edge-cloud kubeconfig create --instance-id xxx

  Create or update a kubeconfig for the Edge Cloud instance with instance ID "xxx" in a custom filepath.
  $ stackit beta edge-cloud kubeconfig create --instance-id xxx --filepath yyy

  Get a kubeconfig for the Edge Cloud instance with instance ID "xxx" without writing it to a file and format the output as json.
  $ stackit beta edge-cloud kubeconfig create --instance-id xxx --disable-writing --output-format json

  Create a kubeconfig for the Edge Cloud instance with instance ID "xxx". This will replace your current kubeconfig file.
  $ stackit beta edge-cloud kubeconfig create --instance-id xxx --overwrite
```

### Options

```
      --disable-writing      Disable writing the kubeconfig to a file.
  -e, --expiration string    Expiration time for the kubeconfig, e.g. 5d. By default, the token is valid for 1h.
  -f, --filepath string      Path to the kubeconfig file. A default is chosen by Kubernetes if not set.
  -h, --help                 Help for "stackit beta edge-cloud kubeconfig create"
      --instance-id string   Edge Cloud instance ID
      --overwrite            Force overwrite the kubeconfig file if it exists.
      --switch-context       Switch to the context in the kubeconfig file to the new context.
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit beta edge-cloud kubeconfig](./stackit_beta_edge-cloud_kubeconfig.md)	 - Provides functionality for Edge Cloud kubeconfig.

