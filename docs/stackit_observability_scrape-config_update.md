## stackit observability scrape-config update

Updates a scrape configuration of an Observability instance

### Synopsis

Updates a scrape configuration of an Observability instance.
The payload can be provided as a JSON string or a file path prefixed with "@".
See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_update for information regarding the payload structure.

```
stackit observability scrape-config update JOB_NAME [flags]
```

### Examples

```
  Update a scrape configuration with name "my-config" from Observability instance "xxx", using an API payload sourced from the file "./payload.json"
  $ stackit observability scrape-config update my-config --payload @./payload.json --instance-id xxx

  Update an scrape configuration with name "my-config" from Observability instance "xxx", using an API payload provided as a JSON string
  $ stackit observability scrape-config update my-config --payload "{...}" --instance-id xxx

  Generate a payload with the current values of a scrape configuration, and adapt it with custom values for the different configuration options
  $ stackit observability scrape-config generate-payload --job-name my-config > ./payload.json
  <Modify payload in file>
  $ stackit observability scrape-configs update my-config --payload @./payload.json
```

### Options

```
  -h, --help                 Help for "stackit observability scrape-config update"
      --instance-id string   Instance ID
      --payload string       Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json
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

