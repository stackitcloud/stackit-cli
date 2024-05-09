## stackit load-balancer observability-credentials

Provides functionality for Load Balancer observability credentials

### Synopsis

Provides functionality for Load Balancer observability credentials. These commands can be used to store and update existing credentials, which are valid to be used for Load Balancer observability. This means, e.g. when using Argus, first of all these credentials must be created for that Argus instance (by using "stackit argus credentials create") and then can be managed for a Load Balancer by using the commands in this group.

```
stackit load-balancer observability-credentials [flags]
```

### Options

```
  -h, --help   Help for "stackit load-balancer observability-credentials"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit load-balancer](./stackit_load-balancer.md)	 - Provides functionality for Load Balancer
* [stackit load-balancer observability-credentials add](./stackit_load-balancer_observability-credentials_add.md)	 - Adds observability credentials to Load Balancer
* [stackit load-balancer observability-credentials cleanup](./stackit_load-balancer_observability-credentials_cleanup.md)	 - Deletes observability credentials unused by any Load Balancer
* [stackit load-balancer observability-credentials delete](./stackit_load-balancer_observability-credentials_delete.md)	 - Deletes observability credentials for Load Balancer
* [stackit load-balancer observability-credentials describe](./stackit_load-balancer_observability-credentials_describe.md)	 - Shows details of observability credentials for Load Balancer
* [stackit load-balancer observability-credentials list](./stackit_load-balancer_observability-credentials_list.md)	 - Lists observability credentials for Load Balancer
* [stackit load-balancer observability-credentials update](./stackit_load-balancer_observability-credentials_update.md)	 - Updates observability credentials for Load Balancer

