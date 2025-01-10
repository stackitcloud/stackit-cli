## stackit observability scrape-config list

Lists all scrape configurations of an Observability instance

### Synopsis

Lists all scrape configurations of an Observability instance.

```
stackit observability scrape-config list [flags]
```

### Examples

```
  List all scrape configurations of Observability instance "xxx"
  $ stackit observability scrape-config list --instance-id xxx

  List all scrape configurations of Observability instance "xxx" in JSON format
  $ stackit observability scrape-config list --instance-id xxx --output-format json

  List up to 10 scrape configurations of Observability instance "xxx"
  $ stackit observability scrape-config list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit observability scrape-config list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
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

* [stackit observability scrape-config](./stackit_observability_scrape-config.md)	 - Provides functionality for scrape configurations in Observability

