## stackit beta public-ip update

Updates a Public IP

### Synopsis

Updates a Public IP.

```
stackit beta public-ip update PUBLIC_IP_ID [flags]
```

### Examples

```
  Update public IP with ID "xxx"
  $ stackit beta public-ip update xxx

  Update public IP with ID "xxx" with new labels
  $ stackit beta public-ip update xxx --labels key=value,foo=bar
```

### Options

```
  -h, --help                    Help for "stackit beta public-ip update"
      --labels stringToString   Labels are key-value string pairs which can be attached to a public IP. E.g. '--labels key1=value1,key2=value2,...' (default [])
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

* [stackit beta public-ip](./stackit_beta_public-ip.md)	 - Provides functionality for public IPs

