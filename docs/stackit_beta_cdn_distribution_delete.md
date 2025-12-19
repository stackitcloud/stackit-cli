## stackit beta cdn distribution delete

Delete a CDN distribution

### Synopsis

Delete a CDN distribution by its ID.

```
stackit beta cdn distribution delete [flags]
```

### Examples

```
  Delete a CDN distribution with ID "xxx"
  $ stackit beta cdn distribution delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta cdn distribution delete"
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

* [stackit beta cdn distribution](./stackit_beta_cdn_distribution.md)	 - Manage CDN distributions

