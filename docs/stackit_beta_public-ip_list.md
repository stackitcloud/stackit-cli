## stackit beta public-ip list

Lists all Public IPs of a project

### Synopsis

Lists all Public IPs of a project.

```
stackit beta public-ip list [flags]
```

### Examples

```
  Lists all public IPs
  $ stackit beta public-ip list

  Lists all public IPs which contains the label xxx
  $ stackit beta public-ip list --label-selector xxx

  Lists all public IPs in JSON format
  $ stackit beta public-ip list --output-format json

  Lists up to 10 public IPs
  $ stackit beta public-ip list --limit 10
```

### Options

```
  -h, --help                    Help for "stackit beta public-ip list"
      --label-selector string   Filter by label
      --limit int               Maximum number of entries to list
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

