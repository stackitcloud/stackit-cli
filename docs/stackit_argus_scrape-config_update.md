## stackit argus scrape-config update

Updates a scrape configuration of an Argus instance

### Synopsis

Updates a scrape configuration of an Argus instance.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_partial_update for information regarding the payload structure.

```
stackit argus scrape-config update JOB_NAME [flags]
```

### Examples

```
  Update a scrape configuration from Argus instance "xxx", using an API payload sourced from the file "./payload.json"
  $ stackit argus scrape-config update my-config --payload @./payload.json --instance-id xxx

  Update an scrape configuration from Argus instance "xxx", using an API payload provided as a JSON string
  $ stackit argus scrape-config update my-config --payload "{...}" --instance-id xxx

  Generate a payload with the current values of a scrape configuration, and adapt it with custom values for the different configuration options
  $ stackit argus scrape-config generate-payload --job-name my-config > ./payload.json
  <Modify payload in file>
  $ stackit argus scrape-configs update my-config --payload @./payload.json
```

### Options

```
  -h, --help                 Help for "stackit argus scrape-config update"
      --instance-id string   Instance ID
      --payload string       Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json
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

