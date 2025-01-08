## stackit beta public-ip describe

Shows details of a Public IP

### Synopsis

Shows details of a Public IP.

```
stackit beta public-ip describe PUBLIC_IP_ID [flags]
```

### Examples

```
  Show details of a public IP with ID "xxx"
  $ stackit beta public-ip describe xxx

  Show details of a public IP with ID "xxx" in JSON format
  $ stackit beta public-ip describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta public-ip describe"
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

