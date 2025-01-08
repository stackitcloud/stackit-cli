## stackit beta public-ip create

Creates a Public IP

### Synopsis

Creates a Public IP.

```
stackit beta public-ip create [flags]
```

### Examples

```
  Create a public IP
  $ stackit beta public-ip create

  Create a public IP with associated resource ID "xxx"
  $ stackit beta public-ip create --associated-resource-id xxx

  Create a public IP with associated resource ID "xxx" and labels
  $ stackit beta public-ip create --associated-resource-id xxx --labels key=value,foo=bar
```

### Options

```
      --associated-resource-id string   Associates the public IP with a network interface or virtual IP (ID)
  -h, --help                            Help for "stackit beta public-ip create"
      --labels stringToString           Labels are key-value string pairs which can be attached to a public IP. E.g. '--labels key1=value1,key2=value2,...' (default [])
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

* [stackit beta public-ip](./stackit_beta_public-ip.md)	 - Provides functionality for public IPs

