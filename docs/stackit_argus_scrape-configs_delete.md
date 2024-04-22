## stackit argus scrape-configs delete

Deletes an Argus Scrape Config

### Synopsis

Deletes an Argus Scrape Config.

```
stackit argus scrape-configs delete JOB_NAME [flags]
```

### Examples

```
  Delete an Argus Scrape config with name "my-config" from Argus instance "xxx"
  $ stackit argus scrape-configs delete my-config --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit argus scrape-configs delete"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit argus scrape-configs](./stackit_argus_scrape-configs.md)	 - Provides functionality for scrape configs in Argus.

