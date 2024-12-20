## stackit beta public-ip delete

Deletes a Public IP

### Synopsis

Deletes a Public IP.
If the public IP is still in use, the deletion will fail


```
stackit beta public-ip delete PUBLIC_IP_ID [flags]
```

### Examples

```
  Delete public IP with ID "xxx"
  $ stackit beta public-ip delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta public-ip delete"
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

