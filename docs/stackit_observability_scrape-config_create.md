## stackit observability scrape-config create

Creates a scrape configuration for an Observability instance

### Synopsis

Creates a scrape configuration job for an Observability instance.
The payload can be provided as a JSON string or a file path prefixed with "@".
If no payload is provided, a default payload will be used.
See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.

```
stackit observability scrape-config create [flags]
```

### Examples

```
  Create a scrape configuration on Observability instance "xxx" using default configuration
  $ stackit observability scrape-config create

  Create a scrape configuration on Observability instance "xxx" using an API payload sourced from the file "./payload.json"
  $ stackit observability scrape-config create --payload @./payload.json --instance-id xxx

  Create a scrape configuration on Observability instance "xxx" using an API payload provided as a JSON string
  $ stackit observability scrape-config create --payload "{...}" --instance-id xxx

  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit observability scrape-config generate-payload > ./payload.json
  <Modify payload in file, if needed>
  $ stackit observability scrape-config create --payload @./payload.json --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit observability scrape-config create"
      --instance-id string   Instance ID
      --payload string       Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit observability scrape-config generate-payload")
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

