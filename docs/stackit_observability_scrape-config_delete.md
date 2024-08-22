## stackit observability scrape-config delete

Deletes a scrape configuration from an Observability instance

### Synopsis

Deletes a scrape configuration from an Observability instance.

```
stackit observability scrape-config delete JOB_NAME [flags]
```

### Examples

```
  Delete a scrape configuration job with name "my-config" from Observability instance "xxx"
  $ stackit observability scrape-config delete my-config --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit observability scrape-config delete"
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

