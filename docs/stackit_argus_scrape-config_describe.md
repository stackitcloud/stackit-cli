## stackit argus scrape-config describe

Shows details of a scrape configuration from an Argus instance

### Synopsis

Shows details of a scrape configuration from an Argus instance.

```
stackit argus scrape-config describe JOB_NAME [flags]
```

### Examples

```
  Get details of a scrape configuration with name "my-config" from Argus instance "xxx"
  $ stackit argus scrape-config describe my-config --instance-id xxx

  Get details of a scrape configuration with name "my-config" from Argus instance "xxx" in a table format
  $ stackit argus scrape-config describe my-config --output-format pretty
```

### Options

```
  -h, --help                 Help for "stackit argus scrape-config describe"
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

* [stackit argus scrape-config](./stackit_argus_scrape-config.md)	 - Provides functionality for scrape configurations in Argus

