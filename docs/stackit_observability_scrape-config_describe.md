## stackit observability scrape-config describe

Shows details of a scrape configuration from an Observability instance

### Synopsis

Shows details of a scrape configuration from an Observability instance.

```
stackit observability scrape-config describe JOB_NAME [flags]
```

### Examples

```
  Get details of a scrape configuration with name "my-config" from Observability instance "xxx"
  $ stackit observability scrape-config describe my-config --instance-id xxx

  Get details of a scrape configuration with name "my-config" from Observability instance "xxx" in JSON format
  $ stackit observability scrape-config describe my-config --output-format json
```

### Options

```
  -h, --help                 Help for "stackit observability scrape-config describe"
      --instance-id string   Instance ID
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

* [stackit observability scrape-config](./stackit_observability_scrape-config.md)	 - Provides functionality for scrape configurations in Observability

